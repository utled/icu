# IDEAS FOR WHERE THIS IS GOING

- Index files on separate machine over ssh
  - Set up access/rsa independently of this program
  - Can the OS lib traverse directory on a remote machine, or do I need to explore a different method for traversal?
  - 

- Feature to copy remote file/dir to local machine
  - Record Inode of origin file
  - Use some form of temp dir for storing
  - Set up monitoring for that temp dir, on change: propose save to remote system
  - Use Rsync for entry transfers
  - Should this program install rsync if not already installed, or be part of the independent setup with ssh access?
  - 

- Feature for entry syncing over multiple machines(?)
  - Is this reasonable to include in this program?
  - Have 1 machine defined as master and the rest as slaves, or be able to sync in multiple directions?
  - 

### TODO
- [] Test traversal on remote machine over ssh