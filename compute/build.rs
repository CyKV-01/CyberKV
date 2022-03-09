fn main() -> Result<(), Box<dyn std::error::Error>> {
    tonic_build::configure().out_dir("src/proto").compile(
        &["../proto/kv.proto", "../proto/coordinator.proto"],
        &["../proto"],
    )?;

    Ok(())
}
