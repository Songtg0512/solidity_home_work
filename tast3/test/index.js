const { ethers,deployments}  = require("hardhat");
const { expect } = require("chai")

describe("start", async function () {
    // it("Should be able to deploy", async function () {
    //     const Contract=await ethers.getContractFactory("NFTAuction")
    //     const contract=await Contract.deploy()
    //     await contract.waitForDeployment()


    //     contract.createAuction(
    //         1000*1000,
    //         ethers.parseEther("0.000000000001"),
    //         ethers.ZeroAddress,
    //         1

    //     )

    //     const auction = await contract.auctions(0)

    //     console.log(auction);


    // });
    let nftAuction;
    let nftAuctionProxy;
    beforeEach(async function () {

        // 部署合约
        await deployments.fixture(["depolyNftAuction"]);
        nftAuctionProxy = await deployments.get("NftAuctionProxy01");

        nftAuction = await ethers.getContractAt(
            "NFTAuction",
            nftAuctionProxy.address
        );  
        await nftAuction.waitForDeployment();
    });
    it("Should be able to deploy", async function () {

        // 1. 检查初始状态
        const initialNextId = await nftAuction.nextAuctionId();
        console.log("初始 nextAuctionId:", initialNextId.toString());
        expect(initialNextId).to.equal(0);
        // 创建拍卖

        nftAuction.createAuction(
            100,
            ethers.parseEther("0.0008"),
            ethers.ZeroAddress,
            0
        )
        const auction = await nftAuction.auctions(0);
        console.log("创建拍卖成功：：", auction);

        // 升级合约
        await deployments.fixture(["upgradeNftAuctionn"]);

        

        const auction2 = await nftAuction.auctions(0);
        console.log("升级后创建拍卖成功：：", auction);
        // expect(auction.startTime).to.equal(auction2.startTime); 




        


    });


   
})

