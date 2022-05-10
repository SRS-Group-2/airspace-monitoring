package aircraft

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._

class ListSimulation extends Simulation {
  val BASE_URL = ""
  val N_USERS = 600

  val httpProtocol = http
    .baseUrl(BASE_URL)
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .doNotTrackHeader("1")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")

  val scn = scenario("ListSimulation")
    .exec(
      http("Aircraft list request")
        .get("/airspace/aircraft/list")
    )
    .pause(5)

  setUp(
    scn.inject(atOnceUsers(N_USERS))
  ).protocols(httpProtocol)
}
