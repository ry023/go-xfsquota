# -*- mode: ruby -*-
# vi: set ft=ruby :
Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/jammy64"
  config.vm.synced_folder "./shared", "/home/vagrant/shared", owner: "vagrant", group: "vagrant"

  config.vm.provider "virtualbox" do |vb|
    disk_file = "./tmp/disk.vdi"
    disk_file_2 = "./tmp/disk2.vdi"
    unless File.exists?(disk_file)
      vb.customize ['createhd', '--filename', disk_file, '--format', 'VDI', '--size', 1024]
    end
    vb.customize ['storageattach', :id, '--storagectl', 'SCSI', '--port', 2, '--device', 0, '--type', 'hdd', '--medium', disk_file]

    unless File.exists?(disk_file_2)
      vb.customize ['createhd', '--filename', disk_file_2, '--format', 'VDI', '--size', 1024]
    end

    vb.customize ['storageattach', :id, '--storagectl', 'SCSI', '--port', 3, '--device', 0, '--type', 'hdd', '--medium', disk_file_2]
  end

  config.vm.provision :shell, inline: $script
end

$script = <<END
mkfs -t xfs /dev/sdc
mkdir /xfs_root
/bin/mount -t xfs -o defaults,uquota,gquota,pquota /dev/sdc /xfs_root
mkfs -t xfs /dev/sdd
mkdir /xfs_root_secondary
/bin/mount -t xfs -o defaults,uquota,gquota,pquota /dev/sdd /xfs_root_secondary
END
