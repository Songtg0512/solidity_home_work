// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract IntToRoman {
    function intToRoman(uint256 num) public pure returns (string memory) {
        require(num > 0 && num < 4000, "Number must be 1-3999");

        // 罗马数字对应表
        uint256[13] memory values = [
            uint256(1000),
            900,
            500,
            400,
            100,
            90,
            50,
            40,
            10,
            9,
            5,
            4,
            1
        ];

        string[13] memory symbols = [
            "M",
            "CM",
            "D",
            "CD",
            "C",
            "XC",
            "L",
            "XL",
            "X",
            "IX",
            "V",
            "IV",
            "I"
        ];

        // 创建一个 bytes 数组用于拼接结果
        bytes memory result;

        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                num -= values[i];
                result = abi.encodePacked(result, symbols[i]);
            }
        }

        return string(result);
    }
}
