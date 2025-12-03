// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MergeSortedArray {
    function merge(
        uint256[] memory a,
        uint256[] memory b
    ) public pure returns (uint256[] memory) {
        uint256 lenA = a.length;
        uint256 lenB = b.length;
        uint256[] memory result = new uint256[](lenA + lenB);

        uint256 i = 0;
        uint256 j = 0;
        uint256 k = 0;

        // 两边同时比较，谁小放谁
        while (i < lenA && j < lenB) {
            if (a[i] <= b[j]) {
                result[k++] = a[i++];
            } else {
                result[k++] = b[j++];
            }
        }

        // 剩余部分直接放入
        while (i < lenA) {
            result[k++] = a[i++];
        }

        while (j < lenB) {
            result[k++] = b[j++];
        }

        return result;
    }
}
