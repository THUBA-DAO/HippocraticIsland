cd circuit_js
node generate_witness.js circuit.wasm ../input.json ../witness.wtns
cd ..
snarkjs plonk setup circuit.r1cs pot14_final.ptau circuit_final.zkey
snarkjs plonk prove circuit_final.zkey witness.wtns proof.json public.json

# export verification key
snarkjs zkey export verificationkey circuit_final.zkey verification_key.json

# following used for verification
snarkjs plonk verify verification_key.json public.json proof.json

# simulate a verification call
snarkjs zkey export soliditycalldata public.json proof.json > proof_hex.txt