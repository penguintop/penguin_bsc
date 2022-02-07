// SPDX-License-Identifier: BSD-3-Clause
pragma solidity ^0.8.0;
pragma abicoder v2;
//import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/utils/math/Math.sol";
import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

//import "@openzeppelin/contracts/presets/ERC20PresetMinterPauser.sol";

interface IERC20 {
    /**
     * @dev Returns the amount of tokens in existence.
     */
    function totalSupply() external view returns (uint256);

    /**
     * @dev Returns the amount of tokens owned by `account`.
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev Moves `amount` tokens from the caller's account to `recipient`.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transfer(address recipient, uint256 amount) external returns (bool);

    /**
     * @dev Returns the remaining number of tokens that `spender` will be
     * allowed to spend on behalf of `owner` through {transferFrom}. This is
     * zero by default.
     *
     * This value changes when {approve} or {transferFrom} are called.
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev Sets `amount` as the allowance of `spender` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * IMPORTANT: Beware that changing an allowance with this method brings the risk
     * that someone may use both the old and the new allowance by unfortunate
     * transaction ordering. One possible solution to mitigate this race
     * condition is to first reduce the spender's allowance to 0 and set the
     * desired value afterwards:
     * https://github.com/ethereum/EIPs/issues/20#issuecomment-263524729
     *
     * Emits an {Approval} event.
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev Moves `amount` tokens from `sender` to `recipient` using the
     * allowance mechanism. `amount` is then deducted from the caller's
     * allowance.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits a {Transfer} event.
     */
    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) external returns (bool);

    /**
     * @dev Emitted when `value` tokens are moved from one account (`from`) to
     * another (`to`).
     *
     * Note that `value` may be zero.
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev Emitted when the allowance of a `spender` for an `owner` is set by
     * a call to {approve}. `value` is the new allowance.
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);
}
interface IERC20Metadata is IERC20 {
    /**
     * @dev Returns the name of the token.
     */
    function name() external view returns (string memory);

    /**
     * @dev Returns the symbol of the token.
     */
    function symbol() external view returns (string memory);

    /**
     * @dev Returns the decimals places of the token.
     */
    function decimals() external view returns (uint8);
}
abstract contract Context {
    function _msgSender() internal view virtual returns (address) {
        return msg.sender;
    }

    function _msgData() internal view virtual returns (bytes calldata) {
        return msg.data;
    }
}

contract ERC20 is Context, IERC20, IERC20Metadata {
    mapping(address => uint256) private _balances;

    mapping(address => mapping(address => uint256)) private _allowances;

    uint256 private _totalSupply;

    string private _name;
    string private _symbol;

    /**
     * @dev Sets the values for {name} and {symbol}.
     *
     * The default value of {decimals} is 18. To select a different value for
     * {decimals} you should overload it.
     *
     * All two of these values are immutable: they can only be set once during
     * construction.
     */
    constructor(string memory name_, string memory symbol_) {
        _name = name_;
        _symbol = symbol_;
    }

    /**
     * @dev Returns the name of the token.
     */
    function name() public view virtual override returns (string memory) {
        return _name;
    }

    /**
     * @dev Returns the symbol of the token, usually a shorter version of the
     * name.
     */
    function symbol() public view virtual override returns (string memory) {
        return _symbol;
    }

    /**
     * @dev Returns the number of decimals used to get its user representation.
     * For example, if `decimals` equals `2`, a balance of `505` tokens should
     * be displayed to a user as `5.05` (`505 / 10 ** 2`).
     *
     * Tokens usually opt for a value of 18, imitating the relationship between
     * Ether and Wei. This is the value {ERC20} uses, unless this function is
     * overridden;
     *
     * NOTE: This information is only used for _display_ purposes: it in
     * no way affects any of the arithmetic of the contract, including
     * {IERC20-balanceOf} and {IERC20-transfer}.
     */
    function decimals() public view virtual override returns (uint8) {
        return 18;
    }

    /**
     * @dev See {IERC20-totalSupply}.
     */
    function totalSupply() public view virtual override returns (uint256) {
        return _totalSupply;
    }

    /**
     * @dev See {IERC20-balanceOf}.
     */
    function balanceOf(address account) public view virtual override returns (uint256) {
        return _balances[account];
    }

    /**
     * @dev See {IERC20-transfer}.
     *
     * Requirements:
     *
     * - `recipient` cannot be the zero address.
     * - the caller must have a balance of at least `amount`.
     */
    function transfer(address recipient, uint256 amount) public virtual override returns (bool) {
        _transfer(_msgSender(), recipient, amount);
        return true;
    }

    /**
     * @dev See {IERC20-allowance}.
     */
    function allowance(address owner, address spender) public view virtual override returns (uint256) {
        return _allowances[owner][spender];
    }

    /**
     * @dev See {IERC20-approve}.
     *
     * NOTE: If `amount` is the maximum `uint256`, the allowance is not updated on
     * `transferFrom`. This is semantically equivalent to an infinite approval.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function approve(address spender, uint256 amount) public virtual override returns (bool) {
        _approve(_msgSender(), spender, amount);
        return true;
    }

    /**
     * @dev See {IERC20-transferFrom}.
     *
     * Emits an {Approval} event indicating the updated allowance. This is not
     * required by the EIP. See the note at the beginning of {ERC20}.
     *
     * NOTE: Does not update the allowance if the current allowance
     * is the maximum `uint256`.
     *
     * Requirements:
     *
     * - `sender` and `recipient` cannot be the zero address.
     * - `sender` must have a balance of at least `amount`.
     * - the caller must have allowance for ``sender``'s tokens of at least
     * `amount`.
     */
    function transferFrom(
        address sender,
        address recipient,
        uint256 amount
    ) public virtual override returns (bool) {
        uint256 currentAllowance = _allowances[sender][_msgSender()];
        if (currentAllowance != type(uint256).max) {
            require(currentAllowance >= amount, "ERC20: transfer amount exceeds allowance");
            unchecked {
                _approve(sender, _msgSender(), currentAllowance - amount);
            }
        }

        _transfer(sender, recipient, amount);

        return true;
    }

    /**
     * @dev Atomically increases the allowance granted to `spender` by the caller.
     *
     * This is an alternative to {approve} that can be used as a mitigation for
     * problems described in {IERC20-approve}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function increaseAllowance(address spender, uint256 addedValue) public virtual returns (bool) {
        _approve(_msgSender(), spender, _allowances[_msgSender()][spender] + addedValue);
        return true;
    }

    /**
     * @dev Atomically decreases the allowance granted to `spender` by the caller.
     *
     * This is an alternative to {approve} that can be used as a mitigation for
     * problems described in {IERC20-approve}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     * - `spender` must have allowance for the caller of at least
     * `subtractedValue`.
     */
    function decreaseAllowance(address spender, uint256 subtractedValue) public virtual returns (bool) {
        uint256 currentAllowance = _allowances[_msgSender()][spender];
        require(currentAllowance >= subtractedValue, "ERC20: decreased allowance below zero");
        unchecked {
            _approve(_msgSender(), spender, currentAllowance - subtractedValue);
        }

        return true;
    }

    /**
     * @dev Moves `amount` of tokens from `sender` to `recipient`.
     *
     * This internal function is equivalent to {transfer}, and can be used to
     * e.g. implement automatic token fees, slashing mechanisms, etc.
     *
     * Emits a {Transfer} event.
     *
     * Requirements:
     *
     * - `sender` cannot be the zero address.
     * - `recipient` cannot be the zero address.
     * - `sender` must have a balance of at least `amount`.
     */
    function _transfer(
        address sender,
        address recipient,
        uint256 amount
    ) internal virtual {
        require(sender != address(0), "ERC20: transfer from the zero address");
        require(recipient != address(0), "ERC20: transfer to the zero address");

        _beforeTokenTransfer(sender, recipient, amount);

        uint256 senderBalance = _balances[sender];
        require(senderBalance >= amount, "ERC20: transfer amount exceeds balance");
        unchecked {
            _balances[sender] = senderBalance - amount;
        }
        _balances[recipient] += amount;

        emit Transfer(sender, recipient, amount);

        _afterTokenTransfer(sender, recipient, amount);
    }

    /** @dev Creates `amount` tokens and assigns them to `account`, increasing
     * the total supply.
     *
     * Emits a {Transfer} event with `from` set to the zero address.
     *
     * Requirements:
     *
     * - `account` cannot be the zero address.
     */
    function _mint(address account, uint256 amount) internal virtual {
        require(account != address(0), "ERC20: mint to the zero address");

        _beforeTokenTransfer(address(0), account, amount);

        _totalSupply += amount;
        _balances[account] += amount;
        emit Transfer(address(0), account, amount);

        _afterTokenTransfer(address(0), account, amount);
    }

    /**
     * @dev Destroys `amount` tokens from `account`, reducing the
     * total supply.
     *
     * Emits a {Transfer} event with `to` set to the zero address.
     *
     * Requirements:
     *
     * - `account` cannot be the zero address.
     * - `account` must have at least `amount` tokens.
     */
    function _burn(address account, uint256 amount) internal virtual {
        require(account != address(0), "ERC20: burn from the zero address");

        _beforeTokenTransfer(account, address(0), amount);

        uint256 accountBalance = _balances[account];
        require(accountBalance >= amount, "ERC20: burn amount exceeds balance");
        unchecked {
            _balances[account] = accountBalance - amount;
        }
        _totalSupply -= amount;

        emit Transfer(account, address(0), amount);

        _afterTokenTransfer(account, address(0), amount);
    }

    /**
     * @dev Sets `amount` as the allowance of `spender` over the `owner` s tokens.
     *
     * This internal function is equivalent to `approve`, and can be used to
     * e.g. set automatic allowances for certain subsystems, etc.
     *
     * Emits an {Approval} event.
     *
     * Requirements:
     *
     * - `owner` cannot be the zero address.
     * - `spender` cannot be the zero address.
     */
    function _approve(
        address owner,
        address spender,
        uint256 amount
    ) internal virtual {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");

        _allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    /**
     * @dev Hook that is called before any transfer of tokens. This includes
     * minting and burning.
     *
     * Calling conditions:
     *
     * - when `from` and `to` are both non-zero, `amount` of ``from``'s tokens
     * will be transferred to `to`.
     * - when `from` is zero, `amount` tokens will be minted for `to`.
     * - when `to` is zero, `amount` of ``from``'s tokens will be burned.
     * - `from` and `to` are never both zero.
     *
     * To learn more about hooks, head to xref:ROOT:extending-contracts.adoc#using-hooks[Using Hooks].
     */
    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual {}

    /**
     * @dev Hook that is called after any transfer of tokens. This includes
     * minting and burning.
     *
     * Calling conditions:
     *
     * - when `from` and `to` are both non-zero, `amount` of ``from``'s tokens
     * has been transferred to `to`.
     * - when `from` is zero, `amount` tokens have been minted for `to`.
     * - when `to` is zero, `amount` of ``from``'s tokens have been burned.
     * - `from` and `to` are never both zero.
     *
     * To learn more about hooks, head to xref:ROOT:extending-contracts.adoc#using-hooks[Using Hooks].
     */
    function _afterTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual {}
}
library SafeMath {
    /**
     * @dev Returns the addition of two unsigned integers, with an overflow flag.
     *
     * _Available since v3.4._
     */
    function tryAdd(uint256 a, uint256 b) internal pure returns (bool, uint256) {
        unchecked {
            uint256 c = a + b;
            if (c < a) return (false, 0);
            return (true, c);
        }
    }

    /**
     * @dev Returns the substraction of two unsigned integers, with an overflow flag.
     *
     * _Available since v3.4._
     */
    function trySub(uint256 a, uint256 b) internal pure returns (bool, uint256) {
        unchecked {
            if (b > a) return (false, 0);
            return (true, a - b);
        }
    }

    /**
     * @dev Returns the multiplication of two unsigned integers, with an overflow flag.
     *
     * _Available since v3.4._
     */
    function tryMul(uint256 a, uint256 b) internal pure returns (bool, uint256) {
        unchecked {
            // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
            // benefit is lost if 'b' is also tested.
            // See: https://github.com/OpenZeppelin/openzeppelin-contracts/pull/522
            if (a == 0) return (true, 0);
            uint256 c = a * b;
            if (c / a != b) return (false, 0);
            return (true, c);
        }
    }

    /**
     * @dev Returns the division of two unsigned integers, with a division by zero flag.
     *
     * _Available since v3.4._
     */
    function tryDiv(uint256 a, uint256 b) internal pure returns (bool, uint256) {
        unchecked {
            if (b == 0) return (false, 0);
            return (true, a / b);
        }
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers, with a division by zero flag.
     *
     * _Available since v3.4._
     */
    function tryMod(uint256 a, uint256 b) internal pure returns (bool, uint256) {
        unchecked {
            if (b == 0) return (false, 0);
            return (true, a % b);
        }
    }

    /**
     * @dev Returns the addition of two unsigned integers, reverting on
     * overflow.
     *
     * Counterpart to Solidity's `+` operator.
     *
     * Requirements:
     *
     * - Addition cannot overflow.
     */
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        return a + b;
    }

    /**
     * @dev Returns the subtraction of two unsigned integers, reverting on
     * overflow (when the result is negative).
     *
     * Counterpart to Solidity's `-` operator.
     *
     * Requirements:
     *
     * - Subtraction cannot overflow.
     */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        return a - b;
    }

    /**
     * @dev Returns the multiplication of two unsigned integers, reverting on
     * overflow.
     *
     * Counterpart to Solidity's `*` operator.
     *
     * Requirements:
     *
     * - Multiplication cannot overflow.
     */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        return a * b;
    }

    /**
     * @dev Returns the integer division of two unsigned integers, reverting on
     * division by zero. The result is rounded towards zero.
     *
     * Counterpart to Solidity's `/` operator.
     *
     * Requirements:
     *
     * - The divisor cannot be zero.
     */
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        return a / b;
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers. (unsigned integer modulo),
     * reverting when dividing by zero.
     *
     * Counterpart to Solidity's `%` operator. This function uses a `revert`
     * opcode (which leaves remaining gas untouched) while Solidity uses an
     * invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     *
     * - The divisor cannot be zero.
     */
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        return a % b;
    }

    /**
     * @dev Returns the subtraction of two unsigned integers, reverting with custom message on
     * overflow (when the result is negative).
     *
     * CAUTION: This function is deprecated because it requires allocating memory for the error
     * message unnecessarily. For custom revert reasons use {trySub}.
     *
     * Counterpart to Solidity's `-` operator.
     *
     * Requirements:
     *
     * - Subtraction cannot overflow.
     */
    function sub(
        uint256 a,
        uint256 b,
        string memory errorMessage
    ) internal pure returns (uint256) {
        unchecked {
            require(b <= a, errorMessage);
            return a - b;
        }
    }

    /**
     * @dev Returns the integer division of two unsigned integers, reverting with custom message on
     * division by zero. The result is rounded towards zero.
     *
     * Counterpart to Solidity's `/` operator. Note: this function uses a
     * `revert` opcode (which leaves remaining gas untouched) while Solidity
     * uses an invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     *
     * - The divisor cannot be zero.
     */
    function div(
        uint256 a,
        uint256 b,
        string memory errorMessage
    ) internal pure returns (uint256) {
        unchecked {
            require(b > 0, errorMessage);
            return a / b;
        }
    }

    /**
     * @dev Returns the remainder of dividing two unsigned integers. (unsigned integer modulo),
     * reverting with custom message when dividing by zero.
     *
     * CAUTION: This function is deprecated because it requires allocating memory for the error
     * message unnecessarily. For custom revert reasons use {tryMod}.
     *
     * Counterpart to Solidity's `%` operator. This function uses a `revert`
     * opcode (which leaves remaining gas untouched) while Solidity uses an
     * invalid opcode to revert (consuming all remaining gas).
     *
     * Requirements:
     *
     * - The divisor cannot be zero.
     */
    function mod(
        uint256 a,
        uint256 b,
        string memory errorMessage
    ) internal pure returns (uint256) {
        unchecked {
            require(b > 0, errorMessage);
            return a % b;
        }
    }
}
/**
@title Chequebook contract without waivers
@author The Swarm Authors
@notice The chequebook contract allows the issuer of the chequebook to send cheques to an unlimited amount of counterparties.
Furthermore, solvency can be guaranteed via hardDeposits
@dev as an issuer, no cheques should be send if the cumulative worth of a cheques send is above the cumulative worth of all deposits
as a beneficiary, we should always take into account the possibility that a cheque bounces (when no hardDeposits are assigned)
*/
contract ERC20SimpleSwap {
  using SafeMath for uint;

  event ChequeCashed(
    address indexed beneficiary,
    address indexed recipient,
    address indexed caller,
    uint totalPayout,
    uint cumulativePayout,
    uint callerPayout
  );
  event ChequeBounced();
  event HardDepositAmountChanged(address indexed beneficiary, uint amount);
  event HardDepositDecreasePrepared(address indexed beneficiary, uint decreaseAmount);
  event HardDepositTimeoutChanged(address indexed beneficiary, uint timeout);
  event Withdraw(uint amount);

  uint public defaultHardDepositTimeout;
  /* structure to keep track of the hard deposits (on-chain guarantee of solvency) per beneficiary*/
  struct HardDeposit {
    uint amount; /* hard deposit amount allocated */
    uint decreaseAmount; /* decreaseAmount substranced from amount when decrease is requested */
    uint timeout; /* issuer has to wait timeout seconds to decrease hardDeposit, 0 implies applying defaultHardDepositTimeout */
    uint canBeDecreasedAt; /* point in time after which harddeposit can be decreased*/
  }

  struct EIP712Domain {
    string name;
    string version;
    uint256 chainId;
  }

  bytes32 public constant EIP712DOMAIN_TYPEHASH = keccak256(
    "EIP712Domain(string name,string version,uint256 chainId)"
  );
  bytes32 public constant CHEQUE_TYPEHASH = keccak256(
    "Cheque(address chequebook,address beneficiary,uint256 cumulativePayout)"
  );
  bytes32 public constant CASHOUT_TYPEHASH = keccak256(
    "Cashout(address chequebook,address sender,uint256 requestPayout,address recipient,uint256 callerPayout)"
  );
  bytes32 public constant CUSTOMDECREASETIMEOUT_TYPEHASH = keccak256(
    "CustomDecreaseTimeout(address chequebook,address beneficiary,uint256 decreaseTimeout)"
  );

  // the EIP712 domain this contract uses
  function domain() internal  returns (EIP712Domain memory) {
    uint256 chainId;
    assembly {
      chainId := chainid()
    }
    return EIP712Domain({
      name: "Chequebook",
      version: "1.0",
      chainId: chainId
    });
  }

  // compute the EIP712 domain separator. this cannot be constant because it depends on chainId
  function domainSeparator(EIP712Domain memory eip712Domain) internal pure returns (bytes32) {
    return keccak256(abi.encode(
        EIP712DOMAIN_TYPEHASH,
        keccak256(bytes(eip712Domain.name)),
        keccak256(bytes(eip712Domain.version)),
        eip712Domain.chainId
    ));
  }

  // recover a signature with the EIP712 signing scheme
  function recoverEIP712(bytes32 hash, bytes memory sig) internal  returns (address) {
    bytes32 digest = keccak256(abi.encodePacked(
        "\x19\x01",
        domainSeparator(domain()),
        hash
    ));
    return ECDSA.recover(digest, sig);
  }

  /* The token against which this chequebook writes cheques */
  ERC20 public token;
  /* associates every beneficiary with how much has been paid out to them */
  mapping (address => uint) public paidOut;
  /* total amount paid out */
  uint public totalPaidOut;
  /* associates every beneficiary with their HardDeposit */
  mapping (address => HardDeposit) public hardDeposits;
  /* sum of all hard deposits */
  uint public totalHardDeposit;
  /* issuer of the contract, set at construction */
  address public issuer;
  /* indicates wether a cheque bounced in the past */
  bool public bounced;

  /**
  @notice sets the issuer, token and the defaultHardDepositTimeout. can only be called once.
  @param _issuer the issuer of cheques from this chequebook (needed as an argument for "Setting up a chequebook as a payment").
  _issuer must be an Externally Owned Account, or it must support calling the function cashCheque
  @param _token the token this chequebook uses
  @param _defaultHardDepositTimeout duration in seconds which by default will be used to reduce hardDeposit allocations
  */
  function init(address _issuer, address _token, uint _defaultHardDepositTimeout) public {
    require(_issuer != address(0), "invalid issuer");
    require(issuer == address(0), "already initialized");
    issuer = _issuer;
    token = ERC20(_token);
    defaultHardDepositTimeout = _defaultHardDepositTimeout;
  }

  /// @return the balance of the chequebook
  function balance() public view returns(uint) {
    return token.balanceOf(address(this));
  }
  /// @return the part of the balance that is not covered by hard deposits
  function liquidBalance() public view returns(uint) {
    return balance().sub(totalHardDeposit);
  }

  /// @return the part of the balance available for a specific beneficiary
  function liquidBalanceFor(address beneficiary) public view returns(uint) {
    return liquidBalance().add(hardDeposits[beneficiary].amount);
  }
  /**
  @dev internal function responsible for checking the issuerSignature, updating hardDeposit balances and doing transfers.
  Called by cashCheque and cashChequeBeneficary
  @param beneficiary the beneficiary to which cheques were assigned. Beneficiary must be an Externally Owned Account
  @param recipient receives the differences between cumulativePayment and what was already paid-out to the beneficiary minus callerPayout
  @param cumulativePayout cumulative amount of cheques assigned to beneficiary
  @param issuerSig if issuer is not the sender, issuer must have given explicit approval on the cumulativePayout to the beneficiary
  */
  function _cashChequeInternal(
    address beneficiary,
    address recipient,
    uint cumulativePayout,
    uint callerPayout,
    bytes memory issuerSig
  ) internal {
    /* The issuer must have given explicit approval to the cumulativePayout, either by being the caller or by signature*/
    if (msg.sender != issuer) {
      require(issuer == recoverEIP712(chequeHash(address(this), beneficiary, cumulativePayout), issuerSig),
      "invalid issuer signature");
    }
    /* the requestPayout is the amount requested for payment processing */
    uint requestPayout = cumulativePayout.sub(paidOut[beneficiary]);
    /* calculates acutal payout */
    uint totalPayout = Math.min(requestPayout, liquidBalanceFor(beneficiary));
    /* calculates hard-deposit usage */
    uint hardDepositUsage = Math.min(totalPayout, hardDeposits[beneficiary].amount);
    require(totalPayout >= callerPayout, "SimpleSwap: cannot pay caller");
    /* if there are some of the hard deposit used, update hardDeposits*/
    if (hardDepositUsage != 0) {
      hardDeposits[beneficiary].amount = hardDeposits[beneficiary].amount.sub(hardDepositUsage);

      totalHardDeposit = totalHardDeposit.sub(hardDepositUsage);
    }
    /* increase the stored paidOut amount to avoid double payout */
    paidOut[beneficiary] = paidOut[beneficiary].add(totalPayout);
    totalPaidOut = totalPaidOut.add(totalPayout);

    /* let the world know that the issuer has over-promised on outstanding cheques */
    if (requestPayout != totalPayout) {
      bounced = true;
      emit ChequeBounced();
    }

    if (callerPayout != 0) {
    /* do a transfer to the caller if specified*/
      require(token.transfer(msg.sender, callerPayout), "transfer failed");
      /* do the actual payment */
      require(token.transfer(recipient, totalPayout.sub(callerPayout)), "transfer failed");
    } else {
      /* do the actual payment */
      require(token.transfer(recipient, totalPayout), "transfer failed");
    }

    emit ChequeCashed(beneficiary, recipient, msg.sender, totalPayout, cumulativePayout, callerPayout);
  }
  /**
  @notice cash a cheque of the beneficiary by a non-beneficiary and reward the sender for doing so with callerPayout
  @dev a beneficiary must be able to generate signatures (be an Externally Owned Account) to make use of this feature
  @param beneficiary the beneficiary to which cheques were assigned. Beneficiary must be an Externally Owned Account
  @param recipient receives the differences between cumulativePayment and what was already paid-out to the beneficiary minus callerPayout
  @param cumulativePayout cumulative amount of cheques assigned to beneficiary
  @param beneficiarySig beneficiary must have given explicit approval for cashing out the cumulativePayout by the sender and sending the callerPayout
  @param issuerSig if issuer is not the sender, issuer must have given explicit approval on the cumulativePayout to the beneficiary
  @param callerPayout when beneficiary does not have ether yet, he can incentivize other people to cash cheques with help of callerPayout
  @param issuerSig if issuer is not the sender, issuer must have given explicit approval on the cumulativePayout to the beneficiary
  */
  function cashCheque(
    address beneficiary,
    address recipient,
    uint cumulativePayout,
    bytes memory beneficiarySig,
    uint256 callerPayout,
    bytes memory issuerSig
  ) public {
    require(
      beneficiary == recoverEIP712(
        cashOutHash(
          address(this),
          msg.sender,
          cumulativePayout,
          recipient,
          callerPayout
        ), beneficiarySig
      ), "invalid beneficiary signature");
    _cashChequeInternal(beneficiary, recipient, cumulativePayout, callerPayout, issuerSig);
  }

  /**
  @notice cash a cheque as beneficiary
  @param recipient receives the differences between cumulativePayment and what was already paid-out to the beneficiary minus callerPayout
  @param cumulativePayout amount requested to pay out
  @param issuerSig issuer must have given explicit approval on the cumulativePayout to the beneficiary
  */
  function cashChequeBeneficiary(address recipient, uint cumulativePayout, bytes memory issuerSig) public {
    _cashChequeInternal(msg.sender, recipient, cumulativePayout, 0, issuerSig);
  }

  /**
  @notice prepare to decrease the hard deposit
  @dev decreasing hardDeposits must be done in two steps to allow beneficiaries to cash any uncashed cheques (and make use of the assgined hard-deposits)
  @param beneficiary beneficiary whose hard deposit should be decreased
  @param decreaseAmount amount that the deposit is supposed to be decreased by
  */
  function prepareDecreaseHardDeposit(address beneficiary, uint decreaseAmount) public {
    require(msg.sender == issuer, "SimpleSwap: not issuer");
    HardDeposit storage hardDeposit = hardDeposits[beneficiary];
    /* cannot decrease it by more than the deposit */
    require(decreaseAmount <= hardDeposit.amount, "hard deposit not sufficient");
    // if hardDeposit.timeout was never set, apply defaultHardDepositTimeout
    uint timeout = hardDeposit.timeout == 0 ? defaultHardDepositTimeout : hardDeposit.timeout;
    hardDeposit.canBeDecreasedAt = block.timestamp + timeout;
    hardDeposit.decreaseAmount = decreaseAmount;
    emit HardDepositDecreasePrepared(beneficiary, decreaseAmount);
  }

  /**
  @notice decrease the hard deposit after waiting the necesary amount of time since prepareDecreaseHardDeposit was called
  @param beneficiary beneficiary whose hard deposit should be decreased
  */
  function decreaseHardDeposit(address beneficiary) public {
    HardDeposit storage hardDeposit = hardDeposits[beneficiary];
    require(block.timestamp >= hardDeposit.canBeDecreasedAt && hardDeposit.canBeDecreasedAt != 0, "deposit not yet timed out");
    /* this throws if decreaseAmount > amount */
    //TODO: if there is a cash-out in between prepareDecreaseHardDeposit and decreaseHardDeposit, decreaseHardDeposit will throw and reducing hard-deposits is impossible.
    hardDeposit.amount = hardDeposit.amount.sub(hardDeposit.decreaseAmount);
    /* reset the canBeDecreasedAt to avoid a double decrease */
    hardDeposit.canBeDecreasedAt = 0;
    /* keep totalDeposit in sync */
    totalHardDeposit = totalHardDeposit.sub(hardDeposit.decreaseAmount);
    emit HardDepositAmountChanged(beneficiary, hardDeposit.amount);
  }

  /**
  @notice increase the hard deposit
  @param beneficiary beneficiary whose hard deposit should be decreased
  @param amount the new hard deposit
  */
  function increaseHardDeposit(address beneficiary, uint amount) public {
    require(msg.sender == issuer, "SimpleSwap: not issuer");
    /* ensure hard deposits don't exceed the global balance */
    require(totalHardDeposit.add(amount) <= balance(), "hard deposit exceeds balance");

    HardDeposit storage hardDeposit = hardDeposits[beneficiary];
    hardDeposit.amount = hardDeposit.amount.add(amount);
    // we don't explicitely set hardDepositTimout, as zero means using defaultHardDepositTimeout
    totalHardDeposit = totalHardDeposit.add(amount);
    /* disable any pending decrease */
    hardDeposit.canBeDecreasedAt = 0;
    emit HardDepositAmountChanged(beneficiary, hardDeposit.amount);
  }

  /**
  @notice allows for setting a custom hardDepositDecreaseTimeout per beneficiary
  @dev this is required when solvency must be guaranteed for a period longer than the defaultHardDepositDecreaseTimeout
  @param beneficiary beneficiary whose hard deposit decreaseTimeout must be changed
  @param hardDepositTimeout new hardDeposit.timeout for beneficiary
  @param beneficiarySig beneficiary must give explicit approval by giving his signature on the new decreaseTimeout
  */
  function setCustomHardDepositTimeout(
    address beneficiary,
    uint hardDepositTimeout,
    bytes memory beneficiarySig
  ) public {
    require(msg.sender == issuer, "not issuer");
    require(
      beneficiary == recoverEIP712(customDecreaseTimeoutHash(address(this), beneficiary, hardDepositTimeout), beneficiarySig),
      "invalid beneficiary signature"
    );
    hardDeposits[beneficiary].timeout = hardDepositTimeout;
    emit HardDepositTimeoutChanged(beneficiary, hardDepositTimeout);
  }

  /// @notice withdraw ether
  /// @param amount amount to withdraw
  // solhint-disable-next-line no-simple-event-func-name
  function withdraw(uint amount) public {
    /* only issuer can do this */
    require(msg.sender == issuer, "not issuer");
    /* ensure we don't take anything from the hard deposit */
    require(amount <= liquidBalance(), "liquidBalance not sufficient");
    require(token.transfer(issuer, amount), "transfer failed");
  }

  function chequeHash(address chequebook, address beneficiary, uint cumulativePayout)
  internal pure returns (bytes32) {
    return keccak256(abi.encode(
      CHEQUE_TYPEHASH,
      chequebook,
      beneficiary,
      cumulativePayout
    ));
  }  

  function cashOutHash(address chequebook, address sender, uint requestPayout, address recipient, uint callerPayout)
  internal pure returns (bytes32) {
    return keccak256(abi.encode(
      CASHOUT_TYPEHASH,
      chequebook,
      sender,
      requestPayout,
      recipient,
      callerPayout
    ));
  }

  function customDecreaseTimeoutHash(address chequebook, address beneficiary, uint decreaseTimeout)
  internal pure returns (bytes32) {
    return keccak256(abi.encode(
      CUSTOMDECREASETIMEOUT_TYPEHASH,
      chequebook,
      beneficiary,
      decreaseTimeout
    ));
  }
}