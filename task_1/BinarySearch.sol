// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    function binarySearch(
        uint256[] memory arr,
        uint256 target
    ) public pure returns (int256) {
        uint256 left = 0;
        uint256 right = arr.length;

        // 注意 right 是开区间（[left, right)）
        while (left < right) {
            uint256 mid = left + (right - left) / 2;

            if (arr[mid] == target) {
                return int256(mid);
            } else if (arr[mid] < target) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }

        return -1; // 没找到
    }
}
