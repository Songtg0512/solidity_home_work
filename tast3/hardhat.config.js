require("@nomicfoundation/hardhat-toolbox");
require("hardhat-deploy");
require("@openzeppelin/hardhat-upgrades");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: "0.8.28",
   namedAccounts: {
    deployer: 0,
    firstAccount: {
      default: 1
    },
    secondAccount: {
      default: 2
    }
  }
};

