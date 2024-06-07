package test

//
//type IDevice2Set interface {
//	SetInt(int)
//}
//type IDevice2Get interface {
//	GetInt() int
//}
//type IDevice2Data interface {
//	IDeviceGet
//	IDeviceSet
//}
//type IDevice2 interface {
//	IDeviceData
//	Marshal() ([]byte, error)
//	Unmarshal([]byte) error
//}
//
//
//type IDeviceSet interface {
//	SetInt(int)
//	SetDevice2(IDevice2Get)
//}
//type IDeviceGet interface {
//	GetInt() int
//	GetDevice2() IDevice2Data
//}
//type IDeviceData interface {
//	IDeviceGet
//	IDeviceSet
//}
//type IDevice interface {
//	IDeviceData
//	message()
//	Marshal() ([]byte, error)
//	Unmarshal([]byte) error
//}
//
//type IUserSet interface {
//	SetAge(int)IUserData
//	SetDevice(IDeviceGet)IUserData
//}
//type IUserGet interface {
//	GetAge() int
//	GetDevice() IDeviceData
//}
//type IUserData interface {
//	IUserGet
//	IUserSet
//}
//type IUser interface {
//	IUserData
//	message()
//	ProtoMessage()
//	Marshal() ([]byte, error)
//	Unmarshal([]byte) error
//}
//
//func (a *Device) marshal() {
//
//	//a.device.Marshal()
//
//}
//
//
//func (a *User) Marshal() {
//
//	MarshalDevice(a.device)
//
//	//a.device.Marshal()
//
//}
//
////////
//
//type marshaler interface {
//	marshal()
//}
//
//func UserDataTOUser(data IUserData)IUser {
//
//}
//
//func MarshalUser(data IUserData) {
//	data.(IUser).Marshal()
//	NewUser()
//	// copy
//	BackUser()
//}
//func MarshalDevice(data IDeviceData) {
//	data.(marshaler).marshal()
//}
//
//func NewDevice() *Device {
//	return nil
//}
//func NewUser(data IUserData) *User {
//	return nil
//}
//func BackUser(IUser)  {
//	return nil
//}
//
//
//type Device struct {}
//
//func (a *Device) GetInt() int {}
//
//type User struct {
//	Device *Device
//	Age int
//	U IUserOneofU
//}
//type IUserOneofU interface {
//	oneof()
//}
//type oneofU1 struct {
//	U1 *U1
//}
//type oneofU2 struct {
//	U2 *U2
//}
//
//func NewOneofU1(u1 *U1) IUserOneofU {
//	return &oneofU1{U1:u1}
//}
//
//func (a *User)SetU1(U1)  {
//
//}
//func (a *User)SetU2(U2)  {
//
//}
//func (a *User)SetU(IUserOneofU)  {
//
//}
//func (a *User)GetU2() U2 {
//	a.U.(oneofU2)
//}
//
//func (o oneofU1) u() {}
//func (o oneofU2) u() {}
//type U1 struct {}
//type U2 struct {}
//
//func (a *User)SetDevice()  {
//	// value IDevice2Get
//	a.Device= value.(*Device)
//
//	NewDevice()
//	// copy
//}
//func (a *User)SetAge() IUserData  {
//}
//
//func sadfsdf()  {
//	var device IDevice = NewIUser()
//
//
//		ConvertToUser(
//	NewUser().SetAge(nil).SetDevice(),
//	).Marshal()
//
//	user := NewUser()
//	user.SetAge().SetDevice()
//	user.Marshal()
//
//	MarshalUser(
//		NewUser().SetAge().SetDevice()
//		)
//
//}
//
//type iMessage interface {
//	Marshal() ([]byte, error)
//	Unmarshal([]byte) error
//}
//
//func Request(data iMessage)  {
//
//}
//
//func RequestToCasino()  {
//	user := NewUser()
//	Request(user)
//}
//func RequestToBalance()  {
//	device := NewDevice()
//	Request(device)
//}
//func ConvertToUser(user IUserData) *User {
//	return nil
//}
