// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./library/safemath.sol";


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



contract StakingHome{
    using SafeMath for uint;
    enum State{
        NOT_INITED,
        COMMON,
        PAUSED,
        STOPPED
    }
    State public state;

    address public admin;
    uint public totalMinerCount;
    uint public totalStakingAmount;
    uint public totalPunishAmount;
    uint public obtainPunishAmount;
    uint public stakingNeedAmount;
    address public tokenAddr;
    struct StakeInfo{
        address lockAddr;
        uint lockStartTime;
        uint lockEndTime;
        uint lockAmount;
        uint256 nodeAddr;
        uint forFeit;
    }
    mapping(address=> StakeInfo) public stakes;
    mapping(uint256=>address) addrMaps;
    uint public defaultLockDuration;

    modifier onlyState(){
        require(state!=State.NOT_INITED,"contract token not inited");
        require(state!=State.PAUSED,"contract paused");
        require(state!=State.STOPPED,"contract stopped");
        _;

    }
    modifier onlyNotInit(){
        require(state ==State.NOT_INITED,"this token contract inited before");
        _;
    }

    modifier onlyAdmin() {
        require(msg.sender == admin, "Only admin");
        _;
    }
    event ChangeStakingNeedAmount(uint);
    event ChangeDefaultLockDuration(uint);
    event NewStaking(StakeInfo);
    event REDEEM(address,uint);
    event PUNISH(address,uint);
    event ObtainPunish(uint);
    event Pause();
    event UnPaused();
    event Stop();


    constructor(){
        state = State.NOT_INITED;
        admin = msg.sender;
        totalMinerCount = 0;
        totalStakingAmount=0;
        stakingNeedAmount = 0;
        totalPunishAmount =0;
        obtainPunishAmount = 0;
        tokenAddr = address(0);
        defaultLockDuration = 30 days;

    }
    function init(address _tokenAddr,uint _stakingNeedAmount,uint _defaultLockDuration) external onlyNotInit onlyAdmin{
        tokenAddr = _tokenAddr;
        stakingNeedAmount = _stakingNeedAmount;
        defaultLockDuration = _defaultLockDuration;
        state = State.COMMON;
    }

    function setStakingNeedAmount(uint amount) onlyAdmin onlyState external{
        stakingNeedAmount = amount;
        emit ChangeStakingNeedAmount(amount);
    }
    
    
    function setDefaultLockDuration(uint lockDuration) onlyAdmin onlyState external{
        defaultLockDuration = lockDuration;
        emit ChangeDefaultLockDuration(lockDuration);
    }

    function Staking(uint256 nodeAddr) external onlyState {
        require(stakes[msg.sender].lockStartTime==0,"address only can staking once");

        bool suc = IERC20(tokenAddr).transferFrom(msg.sender,address(this),stakingNeedAmount);
        require(suc,"Insufficient user balance");
        stakes[msg.sender].lockAddr = msg.sender;
        stakes[msg.sender].lockStartTime = block.timestamp;
        stakes[msg.sender].lockEndTime = block.timestamp.add(defaultLockDuration);
        stakes[msg.sender].lockAmount = stakingNeedAmount;
        stakes[msg.sender].nodeAddr = nodeAddr;
        stakes[msg.sender].forFeit = 0;
        addrMaps[nodeAddr] =msg.sender;
        totalMinerCount = totalMinerCount.add(1);
        totalStakingAmount = totalStakingAmount.add(1);
        emit NewStaking(stakes[msg.sender]);
    }
    function queryStaking(address user) external view returns(StakeInfo memory){
        return stakes[user];
    }

    function queryXwcAddr(uint256 nodeAddr) external view returns(address){
        return addrMaps[nodeAddr];
    }
    function Redeem() external onlyState {
        require(stakes[msg.sender].lockEndTime<block.timestamp,"The assets cannot be redeem until the time of locking up");
        require(stakes[msg.sender].lockAmount>stakes[msg.sender].forFeit,"The assets pledged by the user are all penalized and cannot be redeemed");
        uint canRedeemBalance = stakes[msg.sender].lockAmount.sub(stakes[msg.sender].forFeit);
        require(totalStakingAmount > canRedeemBalance,"unknow error");
        totalStakingAmount = totalStakingAmount- stakes[msg.sender].lockAmount;
        stakes[msg.sender].lockStartTime = 0;
        stakes[msg.sender].lockEndTime = 0;
        stakes[msg.sender].lockAmount = 0;
        stakes[msg.sender].forFeit = 0;
        bool suc = IERC20(tokenAddr).transfer(msg.sender,canRedeemBalance);
        require(suc,"unkown error");
        totalMinerCount = totalMinerCount -1;

        emit REDEEM(msg.sender,canRedeemBalance);

    }
    function ForceRedeem(address user) external onlyAdmin{
        require(stakes[user].lockEndTime>0,"The assets cannot be redeem until the time of locking up");
        require(stakes[user].lockAmount>stakes[user].forFeit,"The assets pledged by the user are all penalized and cannot be redeemed");
        uint canRedeemBalance = stakes[user].lockAmount.sub(stakes[user].forFeit);
        totalStakingAmount = totalStakingAmount- stakes[msg.sender].lockAmount;
        stakes[user].lockStartTime = 0;
        stakes[user].lockEndTime = 0;
        stakes[user].lockAmount = 0;
        stakes[user].forFeit = 0;
        bool suc = IERC20(tokenAddr).transfer(user,canRedeemBalance);
        totalMinerCount = totalMinerCount -1;
        require(suc,"unkown error");
    }

    function Punish(address user,uint amount) external onlyAdmin onlyState{
        require(stakes[user].lockEndTime>0,"The assets cannot be redeem until the time of locking up");
        uint canRedeemBalance = stakes[user].lockAmount.sub(stakes[user].forFeit);
        if (amount > canRedeemBalance){
            totalStakingAmount = totalStakingAmount- stakes[msg.sender].lockAmount;
            totalPunishAmount =totalPunishAmount+ canRedeemBalance;
            stakes[user].lockStartTime = 0;
            stakes[user].lockEndTime = 0;
            stakes[user].lockAmount = 0;
            stakes[user].forFeit = 0;
            totalMinerCount = totalMinerCount -1;
            emit PUNISH(user,canRedeemBalance);
        }else{
            stakes[user].forFeit = stakes[user].forFeit.add(amount);
            totalPunishAmount =totalPunishAmount.add(amount);
            emit PUNISH(user,amount);
        }
        
    }

    function obtainPunish() external onlyAdmin{
        uint amount = totalPunishAmount.sub(obtainPunishAmount);
        obtainPunishAmount = totalPunishAmount;
         bool suc = IERC20(tokenAddr).transfer(msg.sender,amount);
         require(suc,"Insufficient Balance");
         emit ObtainPunish(amount);
    }
    function pause() external onlyAdmin{
        require(state == State.COMMON,"The contract status must be common");
        state = State.PAUSED;
        emit Pause();
    }
    function unPause() external onlyAdmin{
         require(state == State.PAUSED,"The contract status must be paused");
        state = State.COMMON;
        emit UnPaused();
    }
    function stop() external onlyAdmin{
       require(state == State.PAUSED,"The contract status must be paused");
       state = State.STOPPED;
       emit Stop();
    }

}