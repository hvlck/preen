// std

// crates
use clap::{App, Arg, SubCommand};

// local

fn main() {
    let app = App::new("preen")
        .about(env!("CARGO_PKG_DESCRIPTION"))
        .author(env!("CARGO_PKG_AUTHORS"))
        .subcommand(
            SubCommand::with_name("check").about("Check links.").arg(
                Arg::with_name("file")
                    .short("f")
                    .help("Check a given file or directory.")
                    .takes_value(true),
            ),
        )
        .get_matches();

    if let Some(_v) = app.subcommand_matches("check") {}
}
