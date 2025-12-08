// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title BeggingContract - 简单的募捐合约示例
/// @notice 支持 donate、withdraw、getDonation，含 Donation 事件与简单 Top3 排行榜
contract BeggingContract {
    // 合约拥有者（可提现地址）
    address payable public owner;

    // 每个地址累计捐赠金额（单位：wei）
    mapping(address => uint256) private donations;

    // 保存所有唯一捐赠者（用于遍历/排行榜更新）
    address[] private donorList;
    mapping(address => bool) private hasDonated;

    // Donation 事件（可选挑战 1）
    event Donation(address indexed donor, uint256 amount);

    // 简单的 Top-3 排行（可选挑战 2）
    address[3] public topDonors;
    uint256[3] public topDonations;

    // 可选：捐赠时间窗口（如果不需要可将 startTimestamp/endTimestamp 设为 0）
    uint256 public startTimestamp;
    uint256 public endTimestamp;

    // 构造函数：部署者为 owner
    constructor(uint256 _startTimestamp, uint256 _endTimestamp) {
        owner = payable(msg.sender);
        // 设置时间窗口（传 0,0 则表示不限时）
        startTimestamp = _startTimestamp;
        endTimestamp = _endTimestamp;
    }

    // 修饰符：仅 owner 可调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner");
        _;
    }

    // 修饰符：在允许的时间窗口内（如果窗口被设置）
    modifier withinTimeWindow() {
        if (startTimestamp != 0 || endTimestamp != 0) {
            // 如果只设置了 start 或 end 也支持
            if (startTimestamp != 0) {
                require(block.timestamp >= startTimestamp, "Not started");
            }
            if (endTimestamp != 0) {
                require(block.timestamp <= endTimestamp, "Ended");
            }
        }
        _;
    }

    /// @notice 向合约捐赠（必须附带 value）
    /// @dev 使用 payable 接收资金并记录到 mapping
    function donate() external payable withinTimeWindow {
        require(msg.value > 0, "Must send ETH");

        // 如果是第一次捐赠，记录地址到 donorList
        if (!hasDonated[msg.sender]) {
            hasDonated[msg.sender] = true;
            donorList.push(msg.sender);
        }

        // 累加捐赠金额
        donations[msg.sender] += msg.value;

        // 触发事件
        emit Donation(msg.sender, msg.value);

        // 尝试更新 Top3（基于累计捐赠金额）
        _updateTopDonors(msg.sender);
    }

    /// @notice 查询某个地址累计捐赠金额（wei）
    function getDonation(address _addr) external view returns (uint256) {
        return donations[_addr];
    }

    /// @notice 合约所有者提取合约内所有资金（使用 transfer）
    /// @dev 按要求使用 payable 修饰符和 address.transfer
    function withdraw() external payable onlyOwner {
        uint256 bal = address(this).balance;
        require(bal > 0, "No funds");

        // 将所有余额转给 owner（注意：address.transfer 在某些场景可能失败）
        owner.transfer(bal);
    }

    /// @notice 查看所有捐赠者（仅供调试/少量数据使用）
    function getAllDonors() external view returns (address[] memory) {
        return donorList;
    }

    /// @notice 返回当前合约余额（wei）
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }

    /* ============ 内部：Top-3 维护逻辑 ============ */
    /// @dev 在每次 donate 后更新累计榜单（简单实现，O(1) 比较）
    function _updateTopDonors(address donor) internal {
        uint256 total = donations[donor];

        // 若 donor 已在榜内，更新其位置
        for (uint i = 0; i < 3; i++) {
            if (topDonors[i] == donor) {
                topDonations[i] = total;
                // 向上冒泡
                _bubbleUp(i);
                return;
            }
        }

        // donor 不在榜内，判断是否进入榜单
        if (total <= topDonations[2]) return; // 未进入

        // 将 donor 放在末位再冒泡上升
        topDonors[2] = donor;
        topDonations[2] = total;
        _bubbleUp(2);
    }

    /// @dev 冒泡函数：把 index 位置的项向前移动到合适位置
    function _bubbleUp(uint index) internal {
        while (index > 0) {
            if (topDonations[index] > topDonations[index - 1]) {
                // 交换位置
                (topDonors[index], topDonors[index - 1]) = (topDonors[index - 1], topDonors[index]);
                (topDonations[index], topDonations[index - 1]) = (topDonations[index - 1], topDonations[index]);
                index--;
            } else {
                break;
            }
        }
    }

    /* ============ 回退/接收函数 ============ */
    // 允许直接向合约地址发送 ETH（会被记录为 donate）
    receive() external payable {
        // 自动视为 donate（保留 msg.sender）
        if (msg.value > 0) {
            if (!hasDonated[msg.sender]) {
                hasDonated[msg.sender] = true;
                donorList.push(msg.sender);
            }
            donations[msg.sender] += msg.value;
            emit Donation(msg.sender, msg.value);
            _updateTopDonors(msg.sender);
        }
    }

    fallback() external payable {
        // 如果传输数据但还带ETH，也把它当 donate 处理
        if (msg.value > 0) {
            if (!hasDonated[msg.sender]) {
                hasDonated[msg.sender] = true;
                donorList.push(msg.sender);
            }
            donations[msg.sender] += msg.value;
            emit Donation(msg.sender, msg.value);
            _updateTopDonors(msg.sender);
        }
    }
}
