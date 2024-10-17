# example-bridge

Goal: provide a CEX like experience of moving assets between source and destination chains.

Future work / Productionizing:
1. Security
   - intermdiate nodes: 
      - handling of the public / private key pairs. 
      - run on top of some light weight consensus, or use some multi-sig.
   - relayers: 
      - verify state transitions, root hashes, etc, waits for finality.
      - these are trustless entities, we only need one honest relayer since malicous relayers can't fake application state. 
2. Add timeouts to the process
   - e.g. when destination chain fails to include txn within some time period.
3. Add fees to the bridging process
4. Rate limits
   - For withdrawals, limit the amount of token outflow from the destination chain.
5. Add route for withdrawl
   - this should look very similar to the deposit flow. 
6. More rigorous unit and integration testing.
