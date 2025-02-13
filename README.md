
Fallen Miner
==========

![enter image description here](https://i.imgur.com/np4viC8.png)

### Information

Fallen Miner is an open source, silent XMR miner. The readme will detail features, installation and the roadmap for this project.

#### Donate
If you like the product or stuff I push out, feel free to star the project or donate directly:

SOL: F1EtBxf4sPhUsfPdA2jVFfqJ7eLbbkxx4f2ujVhuPrxT
XMR: 85VkL5hw9YceMWVHPGNoFgLxQxw6qwNdF51uAz96WPYmhDYwswVHhoaWPXWjvFGBstGhUNBgR9UvqcqVvYHDmAvcC9yPy4S

Features
-----------------------

 -   Silent XMR Mining  

- SSL support

- Domain fallback support
    
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
**Requirements:**

 - golang (1.18+)
 - GCC - https://dev.to/gamegods3/how-to-install-gcc-in-windows-10-the-easier-way-422j
 
**Client:**
golang is required for both building the client and server, so they will need to be installed.
for building the client, just compile as standard (go build) however if you want a non console version, you will need to declare this in the build arguments. (go build -ldflags="-H=windowsgui")
eg: `go build -ldflags="-H=windowsgui"`

**Server:**
To build the server, some environment variables will need to be changed, specifically CGO_ENABLED. In order to do this open a cmd window and run:

    set CGO_ENABLED=1
You will then be able to build the project.
To compile the server, run `go build` and the server will build, it may take some time on the first build.

to register an account, go to the /register route on the panel and enter your information, then head to /login to login.

Contact me
----------
My hackforums account:  
https://hackforums.net/member.php?action=profile&uid=4880700

## Roadmap

 - Remote config separate from the panel (pastebin etc)
 - Anti Analysis
 - BotKiller

### Issues
If you notice a bug or any other issue that isn't installation, please feel free to contact me directly, or open an issue on github.

