package property

import "math/big"

var (
	PEN_ERC20_PRECISION = big.NewInt(0).Mul(big.NewInt(1_000_000_000), big.NewInt(1_000_000_000))
)

const (
	StakingAdmin   = "0x4Db51b59d2D15BA5cc88A192cAadE49e46c2CCf2"
	StakingAddress = "0xB0725808C12ca6F94874CAFD33Fe33A46291330f"
)
