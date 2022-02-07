package stakingabi

const (
	StakingABIv0_1_0 = `
[
	{
		"inputs": [],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "ChangeDefaultLockDuration",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "ChangeStakingNeedAmount",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"components": [
					{
						"internalType": "address",
						"name": "lockAddr",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "lockStartTime",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "lockEndTime",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "lockAmount",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "nodeAddr",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "forFeit",
						"type": "uint256"
					}
				],
				"indexed": false,
				"internalType": "struct StakingHome.StakeInfo",
				"name": "",
				"type": "tuple"
			}
		],
		"name": "NewStaking",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "ObtainPunish",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "PUNISH",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [],
		"name": "Pause",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "",
				"type": "address"
			},
			{
				"indexed": false,
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"name": "REDEEM",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [],
		"name": "Stop",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [],
		"name": "UnPaused",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "user",
				"type": "address"
			}
		],
		"name": "ForceRedeem",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "user",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "Punish",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "Redeem",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "nodeAddr",
				"type": "uint256"
			}
		],
		"name": "Staking",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "admin",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "defaultLockDuration",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_tokenAddr",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "_stakingNeedAmount",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "_defaultLockDuration",
				"type": "uint256"
			}
		],
		"name": "init",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "obtainPunish",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "obtainPunishAmount",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "pause",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "user",
				"type": "address"
			}
		],
		"name": "queryStaking",
		"outputs": [
			{
				"components": [
					{
						"internalType": "address",
						"name": "lockAddr",
						"type": "address"
					},
					{
						"internalType": "uint256",
						"name": "lockStartTime",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "lockEndTime",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "lockAmount",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "nodeAddr",
						"type": "uint256"
					},
					{
						"internalType": "uint256",
						"name": "forFeit",
						"type": "uint256"
					}
				],
				"internalType": "struct StakingHome.StakeInfo",
				"name": "",
				"type": "tuple"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "nodeAddr",
				"type": "uint256"
			}
		],
		"name": "queryXwcAddr",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "lockDuration",
				"type": "uint256"
			}
		],
		"name": "setDefaultLockDuration",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "amount",
				"type": "uint256"
			}
		],
		"name": "setStakingNeedAmount",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "stakes",
		"outputs": [
			{
				"internalType": "address",
				"name": "lockAddr",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "lockStartTime",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "lockEndTime",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "lockAmount",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "nodeAddr",
				"type": "uint256"
			},
			{
				"internalType": "uint256",
				"name": "forFeit",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "stakingNeedAmount",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "state",
		"outputs": [
			{
				"internalType": "enum StakingHome.State",
				"name": "",
				"type": "uint8"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "stop",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "tokenAddr",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "totalMinerCount",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "totalPunishAmount",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "totalStakingAmount",
		"outputs": [
			{
				"internalType": "uint256",
				"name": "",
				"type": "uint256"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "unPause",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]
`

	StakingDeployedBinv0_1_0 = `0x608060405234801561001057600080fd5b506004361061014d5760003560e01c8063a4a2a9f6116100c3578063d201114a1161007c578063d201114a14610358578063d60e64f014610361578063d95402e714610369578063ebbcaaf614610371578063f7b188a514610384578063f851a4401461038c57600080fd5b8063a4a2a9f6146102d3578063adf8ccaa146102e6578063af5067f5146102f9578063b03081d214610302578063b95d5c2a1461032b578063c19d93fb1461033e57600080fd5b80634d016f35116101155780634d016f35146102685780635840650c146102715780635fbe4d1d14610284578063761aa20c146102af5780638456cb59146102c2578063882d80f8146102ca57600080fd5b8063055608001461015257806307da68f51461016e5780630fd6699b1461017857806316934fc4146101815780632016859714610205575b600080fd5b61015b60055481565b6040519081526020015b60405180910390f35b6101766103a4565b005b61015b60095481565b6101ce61018f3660046115ca565b6007602052600090815260409020805460018201546002830154600384015460048501546005909501546001600160a01b039094169492939192909186565b604080516001600160a01b0390971687526020870195909552938501929092526060840152608083015260a082015260c001610165565b6102186102133660046115ca565b610448565b604051610165919081516001600160a01b031681526020808301519082015260408083015190820152606080830151908201526080808301519082015260a0918201519181019190915260c00190565b61015b60015481565b61017661027f3660046115ca565b6104ef565b600654610297906001600160a01b031681565b6040516001600160a01b039091168152602001610165565b6101766102bd3660046115e5565b6106e2565b6101766107f1565b61015b60045481565b6101766102e13660046115fe565b6108c7565b6101766102f43660046115e5565b61099b565b61015b60035481565b6102976103103660046115e5565b6000908152600860205260409020546001600160a01b031690565b6101766103393660046115e5565b610c9c565b60005461034b9060ff1681565b6040516101659190611647565b61015b60025481565b610176610da4565b610176610ee5565b61017661037f36600461166f565b6111a1565b610176611431565b6000546102979061010090046001600160a01b031681565b60005461010090046001600160a01b031633146103dc5760405162461bcd60e51b81526004016103d390611699565b60405180910390fd5b600260005460ff1660038111156103f5576103f5611631565b146104125760405162461bcd60e51b81526004016103d3906116bd565b6000805460ff191660031781556040517fbedf0f4abfe86d4ffad593d9607fe70e83ea706033d44d24b3b6283cf3fc4f6b9190a1565b61048a6040518060c0016040528060006001600160a01b0316815260200160008152602001600081526020016000815260200160008152602001600081525090565b506001600160a01b03908116600090815260076020908152604091829020825160c0810184528154909416845260018101549184019190915260028101549183019190915260038101546060830152600481015460808301526005015460a082015290565b60005461010090046001600160a01b0316331461051e5760405162461bcd60e51b81526004016103d390611699565b6001600160a01b0381166000908152600760205260409020600201546105565760405162461bcd60e51b81526004016103d3906116ff565b6001600160a01b03811660009081526007602052604090206005810154600390910154116105965760405162461bcd60e51b81526004016103d39061175c565b6001600160a01b038116600090815260076020526040812060058101546003909101546105c2916114cc565b336000908152600760205260409020600301546002549192506105e4916117df565b60029081556001600160a01b0383811660008181526007602052604080822060018101839055948501829055600385018290556005909401819055600654935163a9059cbb60e01b8152600481019290925260248201859052929091169063a9059cbb906044016020604051808303816000875af115801561066a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061068e91906117f6565b90506001805461069e91906117df565b600155806106dd5760405162461bcd60e51b815260206004820152600c60248201526b3ab735b7bbb71032b93937b960a11b60448201526064016103d3565b505050565b60005461010090046001600160a01b031633146107115760405162461bcd60e51b81526004016103d390611699565b6000805460ff16600381111561072957610729611631565b14156107475760405162461bcd60e51b81526004016103d390611818565b600260005460ff16600381111561076057610760611631565b141561077e5760405162461bcd60e51b81526004016103d39061184f565b600360005460ff16600381111561079757610797611631565b14156107b55760405162461bcd60e51b81526004016103d390611878565b60098190556040518181527febe5c19f66b410fe97505092d80e184c1e41bb7707cc42ca705e7f0f7887154f906020015b60405180910390a150565b60005461010090046001600160a01b031633146108205760405162461bcd60e51b81526004016103d390611699565b600160005460ff16600381111561083957610839611631565b146108915760405162461bcd60e51b815260206004820152602260248201527f54686520636f6e747261637420737461747573206d75737420626520636f6d6d60448201526137b760f11b60648201526084016103d3565b6000805460ff191660021781556040517f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff6259190a1565b6000805460ff1660038111156108df576108df611631565b146109365760405162461bcd60e51b815260206004820152602160248201527f7468697320746f6b656e20636f6e747261637420696e69746564206265666f726044820152606560f81b60648201526084016103d3565b60005461010090046001600160a01b031633146109655760405162461bcd60e51b81526004016103d390611699565b600680546001600160a01b0319166001600160a01b0394909416939093179092556005556009556000805460ff19166001179055565b6000805460ff1660038111156109b3576109b3611631565b14156109d15760405162461bcd60e51b81526004016103d390611818565b600260005460ff1660038111156109ea576109ea611631565b1415610a085760405162461bcd60e51b81526004016103d39061184f565b600360005460ff166003811115610a2157610a21611631565b1415610a3f5760405162461bcd60e51b81526004016103d390611878565b3360009081526007602052604090206001015415610a9f5760405162461bcd60e51b815260206004820152601d60248201527f61646472657373206f6e6c792063616e207374616b696e67206f6e636500000060448201526064016103d3565b6006546005546040516323b872dd60e01b815233600482015230602482015260448101919091526000916001600160a01b0316906323b872dd906064016020604051808303816000875af1158015610afb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b1f91906117f6565b905080610b6e5760405162461bcd60e51b815260206004820152601960248201527f496e73756666696369656e7420757365722062616c616e63650000000000000060448201526064016103d3565b33600081815260076020526040902080546001600160a01b0319169091178155426001909101819055600954610ba49190611515565b336000818152600760209081526040808320600281019590955560058054600387015560048601889055909401829055858252600890529190912080546001600160a01b031916909117905560018054610bfd91611515565b6001908155600254610c0e91611515565b600255336000908152600760205260409081902090517f3d9d0020a0d641ffef2f4f9c4c12053788b2362999aa3a1089a67a2dd49fff0891610c909181546001600160a01b031681526001820154602082015260028201546040820152600382015460608201526004820154608082015260059091015460a082015260c00190565b60405180910390a15050565b60005461010090046001600160a01b03163314610ccb5760405162461bcd60e51b81526004016103d390611699565b6000805460ff166003811115610ce357610ce3611631565b1415610d015760405162461bcd60e51b81526004016103d390611818565b600260005460ff166003811115610d1a57610d1a611631565b1415610d385760405162461bcd60e51b81526004016103d39061184f565b600360005460ff166003811115610d5157610d51611631565b1415610d6f5760405162461bcd60e51b81526004016103d390611878565b60058190556040518181527f6f96d542a20cce675aa4c184ef47a85e2bd8e6b5ae3016af332eb4bb8cedea76906020016107e6565b60005461010090046001600160a01b03163314610dd35760405162461bcd60e51b81526004016103d390611699565b6000610dec6004546003546114cc90919063ffffffff16565b600354600490815560065460405163a9059cbb60e01b81523392810192909252602482018390529192506000916001600160a01b03169063a9059cbb906044016020604051808303816000875af1158015610e4b573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e6f91906117f6565b905080610eb55760405162461bcd60e51b8152602060048201526014602482015273496e73756666696369656e742042616c616e636560601b60448201526064016103d3565b6040518281527f38c98a8517e43fafacd84229992dd41467bf97c6a3a9624bf4d06ef554a1b8e790602001610c90565b6000805460ff166003811115610efd57610efd611631565b1415610f1b5760405162461bcd60e51b81526004016103d390611818565b600260005460ff166003811115610f3457610f34611631565b1415610f525760405162461bcd60e51b81526004016103d39061184f565b600360005460ff166003811115610f6b57610f6b611631565b1415610f895760405162461bcd60e51b81526004016103d390611878565b336000908152600760205260409020600201544211610fba5760405162461bcd60e51b81526004016103d3906116ff565b336000908152600760205260409020600581015460039091015411610ff15760405162461bcd60e51b81526004016103d39061175c565b3360009081526007602052604081206005810154600390910154611014916114cc565b905080600254116110565760405162461bcd60e51b815260206004820152600c60248201526b3ab735b737bb9032b93937b960a11b60448201526064016103d3565b3360009081526007602052604090206003015460025461107691906117df565b60029081553360008181526007602052604080822060018101839055938401829055600384018290556005909301819055600654925163a9059cbb60e01b8152600481019290925260248201849052916001600160a01b03169063a9059cbb906044016020604051808303816000875af11580156110f8573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061111c91906117f6565b90508061115a5760405162461bcd60e51b815260206004820152600c60248201526b3ab735b7bbb71032b93937b960a11b60448201526064016103d3565b6001805461116891906117df565b60015560408051338152602081018490527f8b3cbb7f95cd199f1ef5fabeef750deec1a71c5dcfc8045cf51bee4fc68c13019101610c90565b60005461010090046001600160a01b031633146111d05760405162461bcd60e51b81526004016103d390611699565b6000805460ff1660038111156111e8576111e8611631565b14156112065760405162461bcd60e51b81526004016103d390611818565b600260005460ff16600381111561121f5761121f611631565b141561123d5760405162461bcd60e51b81526004016103d39061184f565b600360005460ff16600381111561125657611256611631565b14156112745760405162461bcd60e51b81526004016103d390611878565b6001600160a01b0382166000908152600760205260409020600201546112ac5760405162461bcd60e51b81526004016103d3906116ff565b6001600160a01b038216600090815260076020526040812060058101546003909101546112d8916114cc565b9050808211156113a0573360009081526007602052604090206003015460025461130291906117df565b6002556003546113139082906118a2565b60039081556001600160a01b038416600090815260076020526040812060018082018390556002820183905592810182905560050155805461135591906117df565b600155604080516001600160a01b0385168152602081018390527f5de7354a968c7c4c0c68d2872d7cd6e07f8c2a2569011c47f3f088756bb0a49991015b60405180910390a1505050565b6001600160a01b0383166000908152600760205260409020600501546113c69083611515565b6001600160a01b0384166000908152600760205260409020600501556003546113ef9083611515565b600355604080516001600160a01b0385168152602081018490527f5de7354a968c7c4c0c68d2872d7cd6e07f8c2a2569011c47f3f088756bb0a4999101611393565b60005461010090046001600160a01b031633146114605760405162461bcd60e51b81526004016103d390611699565b600260005460ff16600381111561147957611479611631565b146114965760405162461bcd60e51b81526004016103d3906116bd565b6000805460ff191660011781556040517f472cf038e2a5f33dbaa68760dbf94ab4e159535e6580c0ac63f8202c7c6c0bb29190a1565b600061150e83836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f770000815250611574565b9392505050565b60008061152283856118a2565b90508381101561150e5760405162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f77000000000060448201526064016103d3565b600081848411156115985760405162461bcd60e51b81526004016103d391906118ba565b5060006115a584866117df565b95945050505050565b80356001600160a01b03811681146115c557600080fd5b919050565b6000602082840312156115dc57600080fd5b61150e826115ae565b6000602082840312156115f757600080fd5b5035919050565b60008060006060848603121561161357600080fd5b61161c846115ae565b95602085013595506040909401359392505050565b634e487b7160e01b600052602160045260246000fd5b602081016004831061166957634e487b7160e01b600052602160045260246000fd5b91905290565b6000806040838503121561168257600080fd5b61168b836115ae565b946020939093013593505050565b6020808252600a908201526927b7363c9030b236b4b760b11b604082015260600190565b60208082526022908201527f54686520636f6e747261637420737461747573206d7573742062652070617573604082015261195960f21b606082015260800190565b60208082526038908201527f546865206173736574732063616e6e6f742062652072656465656d20756e746960408201527f6c207468652074696d65206f66206c6f636b696e672075700000000000000000606082015260800190565b60208082526047908201527f5468652061737365747320706c6564676564206279207468652075736572206160408201527f726520616c6c2070656e616c697a656420616e642063616e6e6f742062652072606082015266195919595b595960ca1b608082015260a00190565b634e487b7160e01b600052601160045260246000fd5b6000828210156117f1576117f16117c9565b500390565b60006020828403121561180857600080fd5b8151801515811461150e57600080fd5b60208082526019908201527f636f6e747261637420746f6b656e206e6f7420696e6974656400000000000000604082015260600190565b6020808252600f908201526e18dbdb9d1c9858dd081c185d5cd959608a1b604082015260600190565b60208082526010908201526f18dbdb9d1c9858dd081cdd1bdc1c195960821b604082015260600190565b600082198211156118b5576118b56117c9565b500190565b600060208083528351808285015260005b818110156118e7578581018301518582016040015282016118cb565b818111156118f9576000604083870101525b50601f01601f191692909201604001939250505056fea26469706673582212206022f43ece5e1c5e6151fe002d02b8f836c4f21f4ea94d0b113e2bf7292ab85c64736f6c634300080b0033`
)
