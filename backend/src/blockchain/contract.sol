// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract OrderTracker {
    event OrderUpdateHashStored(uint256 orderId, bytes32 hash);

    mapping(uint256 => bytes32[]) public updateHashes;

    function storeUpdateHash(uint256 orderId, bytes32 hash) public {
        updateHashes[orderId].push(hash);
        emit OrderUpdateHashStored(orderId, hash);
    }

    function getUpdateHash(uint256 orderId) public view returns (bytes32[] memory){
        return updateHashes[orderId];
    }
}