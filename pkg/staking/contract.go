package staking

import (
	"bytes"
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/go-sw3-abi/sw3abi"
	"github.com/penguintop/penguin_bsc/pkg/penguin"
	"github.com/penguintop/penguin_bsc/pkg/property"
	"github.com/penguintop/penguin_bsc/pkg/sctx"
	"github.com/penguintop/penguin_bsc/pkg/stakingabi"
	"github.com/penguintop/penguin_bsc/pkg/transaction"
	"math/big"
	"reflect"
)

var (
	ErrInvalidStaking = errors.New("not a valid staking contract")

	errDecodeABI = errors.New("could not decode abi data")

	stakingABI = transaction.ParseABIUnchecked(stakingabi.StakingABIv0_1_0)

	erc20ABI = transaction.ParseABIUnchecked(sw3abi.ERC20ABIv0_3_1)
)

type Interface interface {
	QueryStaking(ctx context.Context) (bool, error)
	QueryAllowance(ctx context.Context) (*big.Int, error)
	Staking(ctx context.Context) (bool, error)
}

type stakingContract struct {
	owner                  common.Address
	penguinNode            penguin.Address
	stakingContractAddress common.Address
	penTokenAddress        common.Address
	transactionService     transaction.Service
}

func (s *stakingContract) QueryStaking(ctx context.Context) (bool, error) {
	callData, err := stakingABI.Pack("queryStaking", s.owner)
	if err != nil {
		return false, err
	}

	output, err := s.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &s.stakingContractAddress,
		Data: callData,
	})
	if err != nil {
		return false, err
	}

	results, err := stakingABI.Unpack("queryStaking", output)
	if err != nil {
		return false, err
	}

	if len(results) != 1 {
		return false, errDecodeABI
	}

	lockStartTime, ok := abi.ConvertType(reflect.ValueOf(results[0]).Field(1).Interface(), new(big.Int)).(*big.Int)
	if !ok || lockStartTime == nil {
		return false, errDecodeABI
	}

	if lockStartTime.Uint64() == 0 {
		return false, nil
	}
	return true, nil
}

func (s *stakingContract) QueryAllowance(ctx context.Context) (*big.Int, error) {
	callData, err := erc20ABI.Pack("allowance", s.owner, s.stakingContractAddress)
	if err != nil {
		return nil, err
	}

	output, err := s.transactionService.Call(ctx, &transaction.TxRequest{
		To:   &s.penTokenAddress,
		Data: callData,
	})
	if err != nil {
		return nil, err
	}

	results, err := erc20ABI.Unpack("allowance", output)
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, errDecodeABI
	}
	allowance, ok := abi.ConvertType(results[0], new(big.Int)).(*big.Int)
	if !ok || allowance == nil {
		return nil, errDecodeABI
	}
	return allowance, nil
}

func (s *stakingContract) Staking(ctx context.Context) (bool, error) {
	//fmt.Println("QueryAllowance...")
	allowance, err := s.QueryAllowance(ctx)
	if err != nil {
		return false, err
	}
	if allowance.Cmp(big.NewInt(0).Mul(big.NewInt(10000000), property.PEN_ERC20_PRECISION)) < 0 {
		approveAmount := big.NewInt(0).Mul(big.NewInt(50000000), property.PEN_ERC20_PRECISION)
		//fmt.Println("sendApproveTransaction...")
		_, err := s.sendApproveTransaction(ctx, approveAmount)
		if err != nil {
			return false, err
		}
	}
	//fmt.Println("sendStakingTransaction...")
	_, err = s.sendStakingTransaction(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func New(
	owner common.Address,
	penguinNode penguin.Address,
	stakingContractAddress common.Address,
	penTokenAddress common.Address,
	transactionService transaction.Service,
) Interface {
	return &stakingContract{
		owner:                  owner,
		penguinNode:            penguinNode,
		stakingContractAddress: stakingContractAddress,
		penTokenAddress:        penTokenAddress,
		transactionService:     transactionService,
	}
}

func (s *stakingContract) sendApproveTransaction(ctx context.Context, amount *big.Int) (common.Hash, error) {
	callData, err := erc20ABI.Pack("approve", s.stakingContractAddress, big.NewInt(0).Set(amount))
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:       &s.penTokenAddress,
		Data:     callData,
		GasPrice: sctx.GetGasPrice(ctx),
		GasLimit: 0,
		Value:    big.NewInt(0),
	}

	txHash, err := s.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	receipt, err := s.transactionService.WaitForReceipt(ctx, txHash)
	if err != nil {
		return common.Hash{}, err
	}

	if receipt.Status == 0 {
		return common.Hash{}, transaction.ErrTransactionReverted
	}

	return txHash, nil
}

func (s *stakingContract) sendStakingTransaction(ctx context.Context) (common.Hash, error) {
	val := big.NewInt(0).SetBytes(s.penguinNode.Bytes())
	callData, err := stakingABI.Pack("Staking", val)
	if err != nil {
		return common.Hash{}, err
	}

	request := &transaction.TxRequest{
		To:       &s.stakingContractAddress,
		Data:     callData,
		GasPrice: sctx.GetGasPrice(ctx),
		GasLimit: 0,
		Value:    big.NewInt(0),
	}

	txHash, err := s.transactionService.Send(ctx, request)
	if err != nil {
		return common.Hash{}, err
	}

	receipt, err := s.transactionService.WaitForReceipt(ctx, txHash)
	if err != nil {
		return common.Hash{}, err
	}

	if receipt.Status == 0 {
		return common.Hash{}, transaction.ErrTransactionReverted
	}

	return txHash, nil
}

func LookupERC20Address(ctx context.Context, transactionService transaction.Service, stakingContractAddress common.Address) (common.Address, error) {
	callData, err := stakingABI.Pack("tokenAddr")
	if err != nil {
		return common.Address{}, err
	}

	output, err := transactionService.Call(ctx, &transaction.TxRequest{
		To:   &stakingContractAddress,
		Data: callData,
	})
	if err != nil {
		return common.Address{}, err
	}

	results, err := stakingABI.Unpack("tokenAddr", output)
	if err != nil {
		return common.Address{}, err
	}

	if len(results) != 1 {
		return common.Address{}, errDecodeABI
	}

	erc20Address, ok := abi.ConvertType(results[0], new(common.Address)).(*common.Address)
	if !ok || erc20Address == nil {
		return common.Address{}, errDecodeABI
	}
	return *erc20Address, nil
}

func VerifyBytecode(ctx context.Context, backend transaction.Backend, stakingContract common.Address) error {
	code, err := backend.CodeAt(ctx, stakingContract, nil)
	if err != nil {
		return err
	}

	if !bytes.Equal(code, common.FromHex(stakingabi.StakingDeployedBinv0_1_0)) {
		return errors.New("verify byte code, invalid staking contract")
	}

	return nil
}

func VerifyStakingAdmin(ctx context.Context, transactionService transaction.Service, stakingContractAddress common.Address) (bool, error) {
	callData, err := stakingABI.Pack("admin")
	if err != nil {
		return false, err
	}

	output, err := transactionService.Call(ctx, &transaction.TxRequest{
		To:   &stakingContractAddress,
		Data: callData,
	})
	if err != nil {
		return false, err
	}

	results, err := stakingABI.Unpack("admin", output)
	if err != nil {
		return false, err
	}

	if len(results) != 1 {
		return false, errDecodeABI
	}

	erc20Address, ok := abi.ConvertType(results[0], new(common.Address)).(*common.Address)
	if !ok || erc20Address == nil {
		return false, errDecodeABI
	}

	if bytes.Compare(erc20Address.Bytes(), common.FromHex(property.StakingAdmin)) != 0 {
		return false, errors.New("verify staking admin, invalid staking contract admin")
	}

	return true, nil
}
