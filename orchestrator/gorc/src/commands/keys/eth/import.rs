use abscissa_core::{Command, Options, Runnable};

#[derive(Command, Debug, Default, Options)]
pub struct ImportEthKeyCmd {
    #[options(
        short = "n",
        long = "name",
        help = "import private key [name] [privkey]"
    )]
    pub name: String,

    #[options(
        short = "p",
        long = "privkey",
        help = "import private key [name] [privkey]"
    )]
    pub privkey: String,
}

/// The `gorc cosmos import [name] [privkey]` subcommand: import key
impl Runnable for ImportEthKeyCmd {
    fn run(&self) {
        // todo(shella): glue with signatory crate to import key
    }
}