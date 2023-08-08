package iface

/////////////////////////////
////    sql服務介面
/////////////////////////////

type ISqlService interface {
	//sql call stored procedure
	CallSP(string, interface{}, string, ...interface{}) int64
	//sql select,return rowcount
	Select(string, interface{}, string, ...interface{}) int64
	//sql update,return rowsaffected
	Update(string, string, ...interface{}) int64
	//sql delete,return rowsaffected
	Delete(string, string, ...interface{}) int64
	//sql insert,return rowsaffected
	Create(string, string, ...interface{}) int64
	//sql transaction,return rowsaffected
	Transaction(string, []string, ...[]interface{}) int64
}
