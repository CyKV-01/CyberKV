fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure()
        .out_dir("src/proto")
        .type_attribute("node.NodeInfo", "#[derive(::serde::Serialize)]")
        .compile(
            &[
                "../proto/kvs.proto",
                "../proto/coordinator.proto",
                "../proto/node.proto",
            ],
            &["../proto"],
        )?;

    Ok(())
}
