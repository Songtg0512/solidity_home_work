// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Voting {
    mapping(string => uint256) public votes;

    // 记录所有出现过的候选人，便于 reset
    string[] private candidates;
    mapping(string => bool) private exists;

    // 投票给某个地址
    function vote(string memory candidate) external {
        if (!exists[candidate]) {
            // 不存在这个候选人，那么要记录下来
            exists[candidate] = true;
            candidates.push(candidate);
        }
        votes[candidate] += 1;
    }

     // 获取某个候选人的得票数
    function getVotes(string memory candidate) public view returns(uint256) {
        return votes[candidate];
    }

    // 重置所有候选人
    function resetVotes() public {
        for (uint i = 0 ; i < candidates.length; i++) {
            votes[candidates[i]] = 0;
        }
    }
}