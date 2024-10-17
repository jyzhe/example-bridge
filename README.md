# example-bridge

Goal: provide a CEX like experience of moving assets between source and destination chains.

Future work / Productionizing:
1. Security 
   - intermdiate nodes: 
      - handling of the public / private key pairs. 
      - either run on top of some light weight consensus, or use some multi-sig.
   - relayers: 
      - verify state transitions, root hashes, etc. 
      - these are trustless entities, we only need one honest relayer for this to work. 
2. Add timeouts to the process
   - e.g. when destination chain fails to include txn within some time period.
3. Add fees to the bridging process
