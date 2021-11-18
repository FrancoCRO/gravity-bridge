use crate::{
    batch_relaying::relay_batches, find_latest_valset::find_latest_valset,
    logic_call_relaying::relay_logic_calls, valset_relaying::relay_valsets,
};
use ethereum_gravity::utils::{EthClient, get_gravity_id};
use ethers::types::Address as EthAddress;
use gravity_proto::gravity::query_client::QueryClient as GravityQueryClient;
use std::time::{Duration, Instant};
use tokio::time::sleep as delay_for;
use tonic::transport::Channel;

pub const LOOP_SPEED: Duration = Duration::from_secs(17);

/// This function contains the orchestrator primary loop, it is broken out of the main loop so that
/// it can be called in the test runner for easier orchestration of multi-node tests
pub async fn relayer_main_loop(
    eth_client: EthClient,
    grpc_client: GravityQueryClient<Channel>,
    gravity_contract_address: EthAddress,
    eth_gas_price_multiplier: f32,
) {
    let mut grpc_client = grpc_client;
    loop {
        let loop_start = Instant::now();

        let our_ethereum_address = eth_client.address();
        let current_eth_valset =
            find_latest_valset(&mut grpc_client, gravity_contract_address, eth_client.clone()).await;
        if current_eth_valset.is_err() {
            error!("Could not get current valset! {:?}", current_eth_valset);
            continue;
        }
        let current_eth_valset = current_eth_valset.unwrap();

        let gravity_id =
            get_gravity_id(gravity_contract_address, our_ethereum_address, eth_client.clone()).await;
        if gravity_id.is_err() {
            error!("Failed to get GravityID, check your Eth node");
            return;
        }
        let gravity_id = gravity_id.unwrap();

        relay_valsets(
            current_eth_valset.clone(),
            eth_client.clone(),
            &mut grpc_client,
            gravity_contract_address,
            gravity_id.clone(),
            LOOP_SPEED,
        )
        .await;

        relay_batches(
            current_eth_valset.clone(),
            eth_client.clone(),
            &mut grpc_client,
            gravity_contract_address,
            gravity_id.clone(),
            LOOP_SPEED,
            eth_gas_price_multiplier,
        )
        .await;

        relay_logic_calls(
            current_eth_valset,
            eth_client.clone(),
            &mut grpc_client,
            gravity_contract_address,
            gravity_id.clone(),
            LOOP_SPEED,
            eth_gas_price_multiplier
        )
        .await;

        // a bit of logic that tires to keep things running every 5 seconds exactly
        // this is not required for any specific reason. In fact we expect and plan for
        // the timing being off significantly
        let elapsed = Instant::now() - loop_start;
        if elapsed < LOOP_SPEED {
            delay_for(LOOP_SPEED - elapsed).await;
        }
    }
}
