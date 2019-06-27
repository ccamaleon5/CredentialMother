package util

import (
	"encoding/hex"
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	ks "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
)

// PublicKeyBytesToAddress ...
func PublicKeyBytesToAddress(publicKey []byte) common.Address {
	var buf []byte

	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKey[1:]) // remove EC prefix 04
	buf = hash.Sum(nil)
	address := buf[12:]

	return common.HexToAddress(hex.EncodeToString(address))
}

// IsValidAddress validate hex address
func IsValidAddress(iaddress interface{}) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	switch v := iaddress.(type) {
	case string:
		return re.MatchString(v)
	case common.Address:
		return re.MatchString(v.Hex())
	default:
		return false
	}
}

// IsZeroAddress validate if it's a 0 address
func IsZeroAddress(iaddress interface{}) bool {
	var address common.Address
	switch v := iaddress.(type) {
	case string:
		address = common.HexToAddress(v)
	case common.Address:
		address = v
	default:
		return false
	}

	zeroAddressBytes := common.FromHex("0x0000000000000000000000000000000000000000")
	addressBytes := address.Bytes()
	return reflect.DeepEqual(addressBytes, zeroAddressBytes)
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

// CalcGasCost calculate gas cost given gas limit (units) and gas price (wei)
func CalcGasCost(gasLimit uint64, gasPrice *big.Int) *big.Int {
	gasLimitBig := big.NewInt(int64(gasLimit))
	return gasLimitBig.Mul(gasLimitBig, gasPrice)
}

// SigRSV signatures R S V returned as arrays
func SigRSV(isig interface{}) ([32]byte, [32]byte, uint8) {
	var sig []byte
	switch v := isig.(type) {
	case []byte:
		sig = v
	case string:
		sig, _ = hexutil.Decode(v)
	}

	sigstr := common.Bytes2Hex(sig)
	rS := sigstr[0:64]
	sS := sigstr[64:128]
	R := [32]byte{}
	S := [32]byte{}
	copy(R[:], common.FromHex(rS))
	copy(S[:], common.FromHex(sS))
	vStr := sigstr[128:130]
	vI, _ := strconv.Atoi(vStr)
	V := uint8(vI + 27)

	return R, S, V
}

// EqualBytes ...
func EqualBytes(a [32]byte, b [32]byte) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// CreateKeyStore in host directory
func CreateKeyStore(path string, password string) (*ks.KeyStore, string, error){
	keystore := ks.NewKeyStore(path, ks.StandardScryptN, ks.StandardScryptP)
	account, err := keystore.NewAccount(password)
	if err != nil{
		return nil,"",err
	}

	return keystore,account.Address.Hex(),nil
}

//ImportKeyStore from keystore in host directory
func ImportKeyStore(file string, password string) (*ks.KeyStore, *accounts.Account, error){
	jsonBytes, err := ioutil.ReadFile(file)
    if err != nil {
        return nil,nil,err
    }

	ks := ks.NewKeyStore("./keystore", ks.StandardScryptN, ks.StandardScryptP)

    account, err := ks.Import(jsonBytes, password, password)
    if err != nil {
        return nil,nil,err
    }

    if err := os.Remove(file); err != nil {
        return nil,nil,err
	}
	
	return ks,&account,nil
}

// GetKey get a key from KeyStore
func GetKey(file string, password string) (*ks.Key, error){
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ks.DecryptKey(data, password)
	if err != nil {
		return nil, err
	}

	return key, nil
}

// FromECDSA get PrivateKey
func FromECDSA(privateKey *ecdsa.PrivateKey) (string){
	return hex.EncodeToString(crypto.FromECDSA(privateKey))
}
