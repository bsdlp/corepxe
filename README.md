# corepxe

Pronounced "corpse-y", `corepxe` serves to automate the update strategy for
CoreOS hosts that boot using a PXE server.

## Description

[CoreOS](https://coreos.com) is a Linux distribution with the goal of being a
read-only host system for Docker containers. The [update
strategy](https://coreos.com/docs/cluster-management/setup/update-strategies)
is in essence a systemd cronjob that polls the coreos
[omaha](https://github.com/coreos/go-omaha)
[endpoint](https://public.update.core-os.net/v1/update/) to see if there are
any pending updates for your channel, downloads them in the background, and
swaps the updated image into place on reboot. This update strategy is only used
for CoreOS on-disk installations.

CoreOS also supports the ability to do network boot using PXE. This is superior
(in my mind) to installing CoreOS to disk as it allows for easier management of
host system versions in a central location, and you can just make an LV or a
btrfs filesystem and let Docker go crazy.

Currently, updating PXE images is a [manual
process](https://coreos.com/docs/cluster-management/setup/update-strategies/#updating-pxe/ipxe-machines).
This is where `corepxe` comes in. The `update_engine` requests to the coreos
update endpoint is
[proxied](https://coreos.com/docs/cluster-management/setup/update-strategies/#updating-behind-a-proxy)
through `corepxe`, which
[MITMs](https://en.wikipedia.org/wiki/Man-in-the-middle_attack) the response
to extract the URL for the newest image.

When `corepxe` intercepts a request, it will check to see if there's an ongoing
download. If not, then it'll parse the response to see if there is an update
needed.

Snipped example of a response indicating an update is required:
```
  <updatecheck status="ok">
   <urls>
    <url codebase="https://commondatastorage.googleapis.com/update-storage.core-os.net/amd64-usr/410.0.0/"></url>
   </urls>
   <manifest version="410.0.0">
    <packages>
     <package hash="fCMDlzLpTyNnV8++4+kDoqeEuvA=" name="update.gz" size="111882133" required="false"></package>
    </packages>
    <actions>
     <action event="postinstall" ChromeOSVersion="" sha256="MclaAJ7f63k0cHtYs5Wv5dqGuveyXDfbYwDw7X5SaoA=" needsadmin="false" IsDelta="false" DisablePayloadBackoff="true"></action>
    </actions>
   </manifest>
  </updatecheck>
```

If an update is required, then the update is downloaded and placed in the
appropriate directory to be served by
[`tftp`](https://en.wikipedia.org/wiki/Trivial_File_Transfer_Protocol).

If there isn't an update required (`<updatecheck
status="noupdate"></updatecheck>`) then it just passes along that response to
the host that requested a check.

## Notes

```
Example request:

<?xml version="1.0" encoding="UTF-8"?>
<request protocol="3.0" version="CoreOSUpdateEngine-0.1.0.0"
updaterversion="CoreOSUpdateEngine-0.1.0.0" installsource="scheduler"
ismachine="1">
<os version="Chateau" platform="CoreOS" sp="289.0.0"></os>
<app appid="{e96281a6-d1af-4bde-9a0a-97b76e56dc57}" oem="diskless"
version="289.0.0" track="stable" bootid="{fake-client-018}"
machineid="fake-machine-018" lang="en-US" hardware_class="" delta_okay="false"
>
<event eventtype="3" eventresult="2" previousversion=""></event>
</app>
</request>


Example response:

<?xml version="1.0" encoding="UTF-8"?>
<response protocol="3.0" server="update.core-os.net">
 <daystart elapsed_seconds="0"></daystart>
 <app appid="e96281a6-d1af-4bde-9a0a-97b76e56dc57" status="ok">
  <updatecheck status="ok">
   <urls>
    <url codebase="https://commondatastorage.googleapis.com/update-storage.core-os.net/amd64-usr/410.0.0/"></url>
   </urls>
   <manifest version="410.0.0">
    <packages>
     <package hash="fCMDlzLpTyNnV8++4+kDoqeEuvA=" name="update.gz" size="111882133" required="false"></package>
    </packages>
    <actions>
     <action event="postinstall" ChromeOSVersion="" sha256="MclaAJ7f63k0cHtYs5Wv5dqGuveyXDfbYwDw7X5SaoA=" needsadmin="false" IsDelta="false" DisablePayloadBackoff="true"></action>
    </actions>
   </manifest>
  </updatecheck>
 </app>
</response>
```
