# JANP
Just Another Nmap Parser
## Requirements
go get github.com/lair-framework/go-nmap
## Example
```
# ./main -h
    ____.  _____    _______ __________
    |    | /  _  \   \      \\______   \
    |    |/  /_\  \  /   |   \|     ___/
/\__|    /    |    \/    |    \    |    
\________\____|__  /\____|__  /____|    
                 \/         \/          
	Just Another Nmap Parser
Mattia Reggiani - https://github.com/mattiareggiani/JANP - info@mattiareggiani.com

Usage of ./main:
  -alive
    	Print to stout the alive hosts
  -file string
    	Input Nmap XML file
  -format string
    	Format of output file (HTML or CSV)
  -output string
    	Output file (default "output")
  -port int
    	Print to stout the hosts with the port open
  -service string
    	Print to stout the hosts with the service open
      # ./main -file nmapScan.xml -output report -format html
         ____.  _____    _______ __________
         |    | /  _  \   \      \\______   \
         |    |/  /_\  \  /   |   \|     ___/
     /\__|    /    |    \/    |    \    |    
     \________\____|__  /\____|__  /____|    
                      \/         \/          
     	Just Another Nmap Parser
     Mattia Reggiani - https://github.com/mattiareggiani/JANP - info@mattiareggiani.com

     [*] Looks like a valid Nmap XML file
     [*] Parsing...
     [*] Parsed in 323.401µs
     [*] Wrote 91788 bytes
     [+] HTML file successfully created: report.html
```
![alt Screenshot](https://github.com/mattiareggiani/JANP/blob/master/report.JPG)
