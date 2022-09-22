# -*- mode: ruby -*-
# vi: set ft=ruby :
Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/jammy64"

  config.vm.provider "virtualbox" do |vb|
    disk_file = "./tmp/disk.vdi"
    unless File.exists?(disk_file)
      vb.customize ['createhd', '--filename', disk_file, '--format', 'VDI', '--size', 1024]
    end
    vb.customize ['storageattach', :id, '--storagectl', 'SCSI', '--port', 2, '--device', 0, '--type', 'hdd', '--medium', disk_file]
  end

  config.vm.provision :shell, inline: $script
end

$script = <<END
mkfs -t xfs /dev/sdc
mkdir /xfs_root
/bin/mount -t xfs /dev/sdc /xfs_root
END
