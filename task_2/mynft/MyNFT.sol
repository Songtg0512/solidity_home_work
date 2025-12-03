// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// OpenZeppelin ERC721 标准库及扩展
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title MyNFT
/// @notice 完整 ERC721 NFT 合约，支持铸造、批量铸造、销毁、tokenURI
contract MyNFT is ERC721URIStorage, Ownable {
    uint256 private _tokenIds; // NFT 自增 ID

    event NFTMinted(
        address indexed recipient,
        uint256 indexed tokenId,
        string tokenUri,
        uint256 timestamp
    );

    /// @notice 构造函数，设置 NFT 名称和符号
    constructor() ERC721("MyTestNFT", "MTNFT") Ownable(msg.sender) {}

    /// @notice 铸造 NFT
    /// @param recipient NFT 接收地址
    /// @param _tokenURI NFT 元数据链接（IPFS）
    /// @return tokenId 返回铸造的 NFT ID
    function mintNFT(address recipient, string memory _tokenURI) public onlyOwner returns (uint256) {
    _tokenIds++;
    uint256 newItemId = _tokenIds;

    _mint(recipient, newItemId);
    _setTokenURI(newItemId, _tokenURI); // 使用参数 _tokenURI

    return newItemId;
}

    /// @notice 销毁 NFT
    /// @param tokenId 要销毁的 tokenId
    function burnNFT(uint256 tokenId) public {
        require(ownerOf(tokenId) == msg.sender, "Not token owner");
        super._burn(tokenId); // 调用 OpenZeppelin 内部 burn
    }

    /// @notice 查询 NFT 元数据链接
    function tokenURI(
        uint256 tokenId
    ) public view override returns (string memory) {
        return super.tokenURI(tokenId);
    }

    /// @notice ERC165 接口检测
    function supportsInterface(
        bytes4 interfaceId
    ) public view override(ERC721URIStorage) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}
