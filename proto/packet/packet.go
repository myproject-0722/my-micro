package packet

// 消息协议
const (
	CodeSignIn      = 1 // 设备登录
	CodeSignInACK   = 2 // 设备登录回执
	CodeSyncTrigger = 3 // 消息同步触发
	CodeHeadbeat    = 4 // 心跳
	CodeHeadbeatACK = 5 // 心跳回执
	CodeMessage     = 6 // 消息投递
	CodeMessageACK  = 7 // 消息投递回执
)

// Package 消息包
type Package struct {
	Code    int32  // 消息类型
	Content []byte // 消息体
}

const (
	TypeLen       = 4                 // 消息类型字节数组长度
	HeadLen       = 8                 // 消息头部字节数组长度（消息类型字节数组长度+消息长度字节数组长度）
	ContentMaxLen = 4092 * 4          // 消息体最大长度
	BufLen        = ContentMaxLen + 4 // 缓冲buffer字节数组长度
)
