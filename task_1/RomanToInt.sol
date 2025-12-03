// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanToInt {
    function romanToInt(string memory s) public pure returns (uint256) {
        bytes memory b = bytes(s);
        uint256 total = 0;
        uint256 n = b.length;

        for (uint256 i = 0; i < n; i++) {
            uint256 value = charValue(b[i]);
            
            // 如果不是最后一位，并且当前值 < 下一位值，则减
            if (i + 1 < n) {
                uint256 nextValue = charValue(b[i + 1]);
                if (value < nextValue) {
                    total -= value;
                    continue;
                }
            }

            // 否则加
            total += value;
        }

        return total;
    }

    // 返回单个罗马字符对应的数值
    function charValue(bytes1 c) internal pure returns (uint256) {
        if (c == "I") return 1;
        if (c == "V") return 5;
        if (c == "X") return 10;
        if (c == "L") return 50;
        if (c == "C") return 100;
        if (c == "D") return 500;
        if (c == "M") return 1000;
        revert("Invalid Roman numeral");
    }
}
