Cassandra databases contain a special system schema that has tables that hold meta data information on objects such as keyspaces, tables, views, indexes, triggers, and table columns. On newer versions of Cassandra, the name of the system schema holding these tables is system_schema.

Older versions of Cassandra had a system schema named "system" that included similar tables for getting object meta data. This article will be using the system tables defined in version 3 of Cassandra.

Below are example queries that show how to get information on the following Cassandra objects: keyspaces, tables, views, indexes, triggers, and table columns.

Keyspaces
Keyspaces in Cassandra are a similar concept to schemas in databases such as PostgreSQL or Oracle, or databases in databases such as MySQL. Below is an example query for retrieving keyspace information from Cassandra.


select * from system_schema.keyspaces;

Tables
The query below will return information about all tables in a Cassandra database. The keyspace_name column can be used to filter the results by keyspace.


select * from system_schema.tables;

				
Table Columns
The query below will return column information for a table named employee in the sample keyspace.


select * from system_schema.columns where table_name = 'employee' and keyspace_name = 'sample' allow filtering;

Views
The query below will return information about views defined in a Cassandra database.


select * from system_schema.views;

Indexes
The query below will return information about indexes defined in a Cassandra database. The target column includes information about the columns defined in the index.


select * from system_schema.indexes;

Triggers
The query below will return information about triggers defined in a Cassandra database.


select * from system_schema.triggers;



