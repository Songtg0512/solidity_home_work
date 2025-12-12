const { deployments,ethers}=require("hardhat")
const { expect } = require("chai")

describe("Test erc721",async function () {
    it("should be test",async function () {
        await main();
    });
    
});

async function main() {
    const [signer, buyer] = await ethers.getSigners()
    await deployments.fixture(["depolyNftAuction"]);
    const nftAuctionProxy = await deployments.get("NftAuctionProxy01");

    
    

    // 部署erc721
    const Contract=await ethers.getContractFactory("TestERC721")
    const testERC721=await Contract.deploy()
    await testERC721.waitForDeployment()

    const testERC721Address=await testERC721.getAddress()
    console.log("testERC721Address",testERC721Address);

    // minit 10 个nft
    
    for (let i = 0; i < 10; i++) {
        await testERC721.mint(signer.address, i + 1);
    }

    const tokenId = 1
    // 给代理合约授权
    await testERC721.connect(signer).setApprovalForAll(nftAuctionProxy.address, true);

    // 创建拍卖
    const nftAuction=await ethers.getContractAt(
        "NFTAuction",
        nftAuctionProxy.address
    );
    await nftAuction.createAuction(
        10,
        ethers.parseEther("0.01"),
        testERC721Address,
        tokenId    
    );
    const auction=await nftAuction.auctions(0)
    console.log("创建拍卖成功：：", auction);
    // 参与拍卖
    await nftAuction.connect(buyer).placeBid(0,{ value: ethers.parseEther("0.02") });

    // 结束拍卖
    await new Promise((resolve)=>setTimeout(resolve,1000*10))
    await nftAuction.connect(signer).endAuction(0)


    // 验证结果
    const auctionResult=await nftAuction.auctions(0)
    console.log("结束拍卖后读取拍卖成功：：", auctionResult);
    expect(auctionResult.highestBidder).to.equal(buyer.address); 
    expect(auctionResult.hightestPrice).to.equal(ethers.parseEther("0.02")); 

    // 验证nft所有权
    const owner =await testERC721.ownerOf(tokenId);
    console.log("owner::", owner);
    expect(owner).to.equal(buyer.address);

    
    
}