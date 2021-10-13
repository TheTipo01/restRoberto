//! A small synchronous Rust client for the restRoberto TTS service, built on top of ureq.

use std::{io::Read, time::Duration};

/// The HTTP client.
///
/// ```no_run
/// # use restroberto_client::HttpClient;
/// # use std::{io, fs::File, error::Error};
/// fn main() -> Result<(), Box<dyn Error>> {
///     let client = HttpClient::new("https://example.roberto.domain", "a-valid-token");
///     let mut reader = client.get("top left text")?;
///     let mut output = File::create("output.mp3")?;
///     io::copy(&mut reader, &mut output)?;
///     Ok(())
/// }
/// ```
pub struct HttpClient {
    path: String,
    token: String,
    timeout: Duration,
}

impl HttpClient {
    /// Create a new client with a base URL and a token.
    pub fn new(url: &str, token: &str) -> Self {
        Self {
            path: url.to_string(),
            token: token.to_string(),
            timeout: Duration::from_secs(20),
        }
    }

    /// Set a custom timeout.
    pub fn timeout(self, timeout: Duration) -> Self {
        Self { timeout, ..self }
    }

    /// Have the service read the provided text.
    pub fn get(&self, text: &str) -> Result<impl Read + Send, ureq::Error> {
        let query = ureq::get(&(self.path.clone() + "/audio"))
            .query("token", &self.token)
            .query("text", text)
            .timeout(self.timeout);
        let resp = query.call()?;

        if resp.status() != 202 {
            return Err(ureq::Error::Status(resp.status(), resp));
        }

        let req = ureq::get(&(self.path.clone() + "/temp/" + &resp.into_string()? + ".mp3"))
            .timeout(self.timeout);
        let resp = req.call()?;

        if resp.status() != 200 {
            return Err(ureq::Error::Status(resp.status(), resp));
        }

        Ok(resp.into_reader())
    }
}
