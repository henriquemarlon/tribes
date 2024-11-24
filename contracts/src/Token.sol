// SPDX-License-Identifier: MIT
// Compatible with OpenZeppelin Contracts ^5.0.0
pragma solidity ^0.8.22;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract Token is ERC20 {
    constructor()
        ERC20("USDT", "USDT")
    {}

    function mint(address to, uint256 amount) public {
        _mint(to, amount);
    }
}