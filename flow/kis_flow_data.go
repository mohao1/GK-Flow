package flow

import (
	"context"
	"errors"
	"fmt"
	"github.com/patrickmn/go-cache"
	"kis-folw/common"
	"kis-folw/config"
	"kis-folw/log"
	"kis-folw/metrics"
	"time"
)

// CommitRow 提交Flow数据, 一行数据，如果是批量数据可以提交多次
func (flow *KisFlow) CommitRow(row interface{}) error {

	flow.buffer = append(flow.buffer, row)

	return nil
}

// commitSrcData 提交当前Flow的数据源数据, 表示首次提交当前Flow的原始数据源
// 将flow的临时数据buffer，提交到flow的data中,(data为各个Function层级的源数据备份)
// 会清空之前所有的flow数据
func (flow *KisFlow) commitSrcData(ctx context.Context) error {

	// 制作批量数据batch
	dataCnt := len(flow.buffer)
	batch := make(common.KisRowArr, 0, dataCnt)

	for _, row := range flow.buffer {
		batch = append(batch, row)
	}

	// 清空之前所有数据
	flow.clearData(flow.data)

	// 首次提交，记录flow原始数据
	// 因为首次提交，所以PrevFunctionId为FirstVirtual 因为没有上一层Function
	flow.data[common.FunctionIdFirstVirtual] = batch

	// 清空缓冲Buf
	flow.buffer = flow.buffer[0:0]

	// 首次提交数据源数据，进行统计数据总量
	if config.GlobalConfig.EnableProm == true {
		// 统计数据总量 Metrics.DataTota 指标累计加1
		metrics.Metrics.DataTotal.Add(float64(dataCnt))
		//统计当前Flow数量指标
		metrics.Metrics.FlowDataTotal.WithLabelValues(flow.Name).Add(float64(dataCnt))
	}

	log.Logger().DebugFX(ctx, "====> After CommitSrcData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

// clearData 清空flow所有数据
func (flow *KisFlow) clearData(data common.KisDataMap) {
	for k := range data {
		delete(data, k)
	}
}

// commitCurData 提交Flow当前执行Function的结果数据
func (flow *KisFlow) commitCurData(ctx context.Context) error {

	//判断本层计算是否有结果数据,如果没有则退出本次Flow Run循环
	if len(flow.buffer) == 0 {
		flow.abort = true
		return nil
	}

	// 制作批量数据batch
	batch := make(common.KisRowArr, 0, len(flow.buffer))

	//如果strBuf为空，则没有添加任何数据
	for _, row := range flow.buffer {
		batch = append(batch, row)
	}

	//将本层计算的缓冲数据提交到本层结果数据中
	flow.data[flow.ThisFunctionId] = batch

	//fmt.Println(flow.data[flow.ThisFunctionId])
	//fmt.Println(flow.ThisFunctionId)
	//fmt.Println(flow.PrevFunctionId)

	//清空缓冲Buf
	flow.buffer = flow.buffer[0:0]

	log.Logger().DebugFX(ctx, " ====> After commitCurData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

// getCurData 获取flow当前Function层级的输入数据
func (flow *KisFlow) getCurData() (common.KisRowArr, error) {
	if flow.PrevFunctionId == "" {
		return nil, errors.New(fmt.Sprintf("flow.PrevFunctionId is not set"))
	}

	if _, ok := flow.data[flow.PrevFunctionId]; !ok {
		return nil, errors.New(fmt.Sprintf("[%s] is not in flow.data", flow.PrevFunctionId))
	}

	return flow.data[flow.PrevFunctionId], nil
}

// commitReuseData
func (flow *KisFlow) commitReuseData(ctx context.Context) error {
	// 判断上层是否有结果数据, 如果没有则退出本次Flow Run循环
	if len(flow.data[flow.PrevFunctionId]) == 0 {
		flow.abort = true
		return nil
	}

	// 本层结果数据等于上层结果数据(复用上层结果数据到本层)
	flow.data[flow.ThisFunctionId] = flow.data[flow.PrevFunctionId]

	// 清空缓冲Buf (如果是ReuseData选项，那么提交的全部数据，都将不会携带到下一层)
	flow.buffer = flow.buffer[0:0]

	log.Logger().DebugFX(ctx, " ====> After commitReuseData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

// commitVoidData
func (flow *KisFlow) commitVoidData(ctx context.Context) error {
	if len(flow.buffer) != 0 {
		return nil
	}

	// 制作空数据
	batch := make(common.KisRowArr, 0)

	// 将本层计算的缓冲数据提交到本层结果数据中
	flow.data[flow.ThisFunctionId] = batch

	log.Logger().DebugFX(ctx, " ====> After commitVoidData, flow_name = %s, flow_id = %s\nAll Level Data =\n %+v\n", flow.Name, flow.Id, flow.data)

	return nil
}

func (flow *KisFlow) GetCacheData(key string) interface{} {

	if data, found := flow.cache.Get(key); found {
		return data
	}
	return nil
}

func (flow *KisFlow) SetCacheData(key string, value interface{}, Exp time.Duration) {
	if Exp == common.DefaultExpiration {
		flow.cache.Set(key, value, cache.DefaultExpiration)
	} else {
		flow.cache.Set(key, value, Exp)
	}
}

// GetMetaData 得到当前Flow对象的临时数据
func (flow *KisFlow) GetMetaData(key string) interface{} {
	flow.mLock.RLock()
	defer flow.mLock.RUnlock()

	if data, ok := flow.metaData[key]; !ok {
		return nil
	} else {
		return data
	}

}

// SetMetaData 设置当前Flow对象的临时数据
func (flow *KisFlow) SetMetaData(key string, value interface{}) {
	flow.mLock.Lock()
	defer flow.mLock.RLock()

	flow.metaData[key] = value
}

// GetFuncParam 得到Flow的当前正在执行的Function的配置默认参数，取出一对key-value
func (flow *KisFlow) GetFuncParam(key string) string {
	flow.fplock.RLock()
	defer flow.fplock.Unlock()

	if param, ok := flow.funcParams[flow.ThisFunctionId]; ok {
		if value, ok := param[key]; ok {
			return value
		}
	}

	return ""
}

// GetFuncParamAll 得到Flow的当前正在执行的Function的配置默认参数，取出全部Key-Value
func (flow *KisFlow) GetFuncParamAll() config.FParam {
	flow.fplock.RLock()
	defer flow.fplock.RUnlock()

	if param, ok := flow.funcParams[flow.ThisFunctionId]; !ok {
		return nil
	} else {
		return param
	}
}

// GetFuncParamsAllFuncs 得到Flow中所有Function的FuncParams，取出全部Key-Value
func (flow *KisFlow) GetFuncParamsAllFuncs() map[string]config.FParam {
	flow.fplock.RLock()
	defer flow.fplock.RUnlock()

	return flow.funcParams
}
