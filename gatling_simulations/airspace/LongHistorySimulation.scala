package aircraft

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._

class LongHistorySimulation extends Simulation {
  val BASE_URL = System.getProperty("base_url")
  val N_USERS = Integer.getInteger("users", 20)
  val RAMP = java.lang.Long.getLong("ramp", 0)
  
  val toTimestamp: Long = System.currentTimeMillis / 1000
  val fromTimestamp: Long = toTimestamp - 30 * 24 * 60 * 60 // one month

  val httpProtocol = http
    .baseUrl(BASE_URL)
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .doNotTrackHeader("1")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")

  val scn = scenario("LongHistorySimulation")
    .during(RAMP) {
      exec(
        http("Aircraft list request")
          .get(s"/airspace/history?resolution=hour&to=${toTimestamp}&from=${fromTimestamp}")
      )
    }
    .pause(5)

  setUp(
    scn.inject(atOnceUsers(N_USERS))
  ).protocols(httpProtocol)
}
