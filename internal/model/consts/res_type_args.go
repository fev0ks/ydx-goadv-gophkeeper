package consts

import "ydx-goadv-gophkeeper/internal/model/enum"

const (
	LoginPasswordArg = "lp"
	FileArg          = "fl"
	BankCardArg      = "bc"
)

var (
	ArgToType = map[string]enum.ResourceType{
		LoginPasswordArg: enum.LoginPassword,
		FileArg:          enum.File,
		BankCardArg:      enum.BankCard,
	}

	TypeToArg = map[enum.ResourceType]string{
		enum.LoginPassword: LoginPasswordArg,
		enum.File:          FileArg,
		enum.BankCard:      BankCardArg,
	}
)
