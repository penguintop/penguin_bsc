package factoryabi

const (
	FactoryABIv0_1_0 = `
[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "issuer",
				"type": "address"
			},
			{
				"internalType": "uint256",
				"name": "defaultHardDepositTimeoutDuration",
				"type": "uint256"
			},
			{
				"internalType": "bytes32",
				"name": "salt",
				"type": "bytes32"
			}
		],
		"name": "deploySimpleSwap",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_ERC20Address",
				"type": "address"
			}
		],
		"stateMutability": "nonpayable",
		"type": "constructor"
	},
	{
		"anonymous": false,
		"inputs": [
			{
				"indexed": false,
				"internalType": "address",
				"name": "contractAddress",
				"type": "address"
			}
		],
		"name": "SimpleSwapDeployed",
		"type": "event"
	},
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"name": "deployedContracts",
		"outputs": [
			{
				"internalType": "bool",
				"name": "",
				"type": "bool"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "ERC20Address",
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
		"name": "master",
		"outputs": [
			{
				"internalType": "address",
				"name": "",
				"type": "address"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]
`

	FactoryDeployedBinv_0_1_0 = "0x608060405234801561001057600080fd5b506004361061004c5760003560e01c806315efd8a714610051578063a6021ace14610081578063c70242ad14610094578063ee97f7f3146100c7575b600080fd5b61006461005f3660046102b1565b6100da565b6040516001600160a01b0390911681526020015b60405180910390f35b600154610064906001600160a01b031681565b6100b76100a23660046102e4565b60006020819052908152604090205460ff1681565b6040519015158152602001610078565b600254610064906001600160a01b031681565b60025460408051336020820152908101839052600091829161011e916001600160a01b031690606001604051602081830303815290604052805190602001206101eb565b6001546040516343431f6360e11b81526001600160a01b0388811660048301529182166024820152604481018790529192508216906386863ec690606401600060405180830381600087803b15801561017657600080fd5b505af115801561018a573d6000803e3d6000fd5b505050506001600160a01b03811660008181526020818152604091829020805460ff1916600117905590519182527fc0ffc525a1c7689549d7f79b49eca900e61ac49b43d977f680bcc3b36224c004910160405180910390a1949350505050565b6000604051733d602d80600a3d3981f3363d3d373d3d3d363d7360601b81528360601b60148201526e5af43d82803e903d91602b57fd5bf360881b6028820152826037826000f59150506001600160a01b03811661028f5760405162461bcd60e51b815260206004820152601760248201527f455243313136373a2063726561746532206661696c6564000000000000000000604482015260640160405180910390fd5b92915050565b80356001600160a01b03811681146102ac57600080fd5b919050565b6000806000606084860312156102c657600080fd5b6102cf84610295565b95602085013595506040909401359392505050565b6000602082840312156102f657600080fd5b6102ff82610295565b939250505056fea264697066735822122040610c4f98e37772fe0758efb939a51ffdcb7b61e1828643591078e2064350cc64736f6c634300080b0033"
)
