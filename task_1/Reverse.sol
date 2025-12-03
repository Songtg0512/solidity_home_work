// SPDX-License-Identifier: MIT
pragma solidity ^0.8.5;

contract ReverseString {
    // 反转字符串
    function reverse(string memory s) public pure returns (string memory) {
        bytes memory b = bytes(s);
        uint len = b.length;

        for (uint i = 0; i < len / 2; i++) {
            // 做交换
            bytes1 temp = b[i];
            b[i] = b[len - 1 - i];
            b[len - 1 - i] = temp;
        }

        return string(b);
    }
}
