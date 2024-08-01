#!/bin/bash

# Export the ledger log
 curl -X POST http://localhost:8080/v2/default/logs/export -o /ledger.export.log