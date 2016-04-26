## NAB2016: For issue tracking join the waffle.io kanban board:
https://waffle.io/c4labs/c4/join

## C4 - The Cinema Content Creation Cloud
C4 is a framework for coordinating automation in media production across many physical locations for the full life cycle of a production.

C4 features the ability to unambiguously track assets through different organizations even if those organizations do not use the framework internally. All files can be identified constantly *in flight*, meaning that regardless of where or when a file is encountered it always has the same identification.  This is done without the need for any special handling such as standardized naming conventions, or centralized coordinating service.  Computers don't even need to be connected to the Internet to accurately identify files.

Leveraging this universal identification system the framework establishes the notion of "**indelible metadata**", or in other words information about assets that is permanently bound to the asset. The indelible metadata can be described for automated systems as well as human users via the C4 domain language **c4lang**.

C4lang describes relationships between assets declaratively in an easy to read format that is also easy for software to parse based on YAML.  The language makes it easy to describe the dependency graph of assets, metadata, and processes that go into describing a particular result.

C4lang also promotes an iterative style of construction nostalgic of node based compositing tools, in which one progressively refines the workflow as more detailed requirements emerge.  

The framework has 4 major components to address the key difficulties of production in the cloud.

1. C4id for file identification.

2. C4lang for metadata and workflow descriptions.

3. Simplified standards for API implementations that enable interoperability.

4. A robust, strong cryptography base, security model.

C4 is endorsed by ETC the Entertainment Technology Center, major studios, industry software vendors, and is rapidly becoming the de facto standard for cloud production.

## Tools
C4 consistent of fee open source command line tool and demon.  Users interact with the active c4d runtime via the command line tool.

# Framework Modules
The C4 framework is broken down into the following modules.

- Assets
  + C4id the C4 data identification system.
- Language
  + C4lang declarative dependency graph language.
- Entity
  + User and organizational identification and access management.
- FS
  + C4fs Filesystem and file system related tools.
- Net
  + Data transmission, and network services.
- Process
  + Mechanisms for process definition, execution, and compute resource management.
- Storage
  + Database and file storage services.
- Metrics
  + Defines a cost modeling mechanism for billing and budgeting.



*(More to come once we have a clear idea of what will be included in the NAB software release)*



--- 

# Older Readme (needs cleanup)

## c4 command line

The `c4` command line tool in `cmd/c4` is a 'client' application that is the reference implementation of a c4 framework front end.  

##### Features

- Generates  c4ids
- Recursively traverse a file system and identify all files and folders
- Output basic metadata for identified files, folders, and links

##### Road map
The following list of features are planned, but have not yet been prioritized 

- C4lang parsing/generation
- Server interaction for local or remote persistent directives
- Multi target file copy
- Remote file copy
- File system verification
- Entity identity generation
- Code signing


## c4d
The `c4` daemon is a server that provides the c4 execution environment.  It can provide services for local users and remote connections.  Similar to ssh *All* communication with c4d is encrypted, and requires that the sender is identified (via PKI) and approved to execute code on this server.

`c4d` provides services such as watching folders, file transfer, connection management, and resource management.  These services are evoked using c4lang typicall send by the c4 command line tool or other client.

#### Features

#### Road map

