package main

import (
	fshttp "golearn/freeswitch/fseslhttp"
	fstoesl "golearn/freeswitch/fstoesl"
)

func main() {
	var sccfsinfo fstoesl.Fseslinfo
	sccfsinfo.Fseslclientrun()
	fshttp.Setsccfsinfo(&sccfsinfo)
	select {}

}
