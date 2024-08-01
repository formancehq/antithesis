#!/bin/bash

# Export the ledger log
 curl -X POST http://localhost:3068/v2/xyz/logs/export -o /ledger.export.log