package week02

import (
	"github.com/pkg/errors"
	"fmt"
	"log"
)

/*
问题：1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。
为什么，应该怎么做请写出代码？

答：应该wrap，根据规范，作为第三方库的sql操作库不应该使用 pkg/errors 包装错误，避免重复包装。
	作为服务的数据库操作则需要 wrap 错误，用于追踪该error的调用栈信息
 */

func mysqlFunc()(interface{}, error){
	return nil, fmt.Errorf("sql.ErrNoRows")
}

func dao()(interface{}, error){
	result, err := mysqlFunc()
	if err!=nil{
		return nil, errors.Wrap(err)
	}
	return result,nil
}

func service() interface{} {
	message,err := dao()
	if err!=nil{
		message = fmt.Sprintf("operation failed: %+v", err)
		log.Println(message)
	}
	return message
}