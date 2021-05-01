package utils

import (
	"syscall"
	"unsafe"
)

// DATA_BLOB https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/aa381414(v=vs.85)
type DATA_BLOB struct {
	cbData uint32
	pbData *byte
}

func NewBlob(d []byte) *DATA_BLOB {
	if len(d) == 0 {
		return new(DATA_BLOB)
	}
	return &DATA_BLOB{
		pbData: &d[0],
		cbData: uint32(len(d)),
	}
}

func (dataBlob *DATA_BLOB) ToByteArray() []byte {
	d := make([]byte, dataBlob.cbData)
	copy(d, (*[1 << 30]byte)(unsafe.Pointer(dataBlob.pbData))[:])
	return d
}

// free pDataOut memorize
// kernel32 := syscall.NewLazyDLL("Kernel32.dll")
// procLocalFree := kernel32.NewProc("LocalFree")
// defer procLocalFree.Call(uintptr(unsafe.Pointer(pDataOut.pbData)))

// Go call dll reference https://medium.com/@justen.walker/breaking-all-the-rules-using-go-to-call-windows-api-2cbfd8c79724
// CryptUnprotectData  https://docs.microsoft.com/en-us/windows/win32/api/dpapi/nf-dpapi-cryptunprotectdata
func CryptUnprotectData(pDataIn []byte) ([]byte, error) {
	dllCrypt32 := syscall.NewLazyDLL("Crypt32.dll")
	procDecryptData := dllCrypt32.NewProc("CryptUnprotectData")

	var pDataOut DATA_BLOB
	r, _, err := procDecryptData.Call(uintptr(unsafe.Pointer(NewBlob(pDataIn))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&pDataOut)))
	if r == 0 {
		return nil, err
	}

	return pDataOut.ToByteArray(), nil
}
