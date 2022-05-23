package aircraft

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._

class ListPositionSimulation extends Simulation {
  val BASE_URL = System.getProperty("base_url")
  val N_USERS = Integer.getInteger("users", 20)
  val RAMP = java.lang.Long.getLong("ramp", 0)
  val WS_CONNECTION_TIME = 100
  val EXPLORE_API_USER_TIME = 60

  val httpProtocol = http
    .baseUrl(BASE_URL)
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .doNotTrackHeader("1")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")

  val scn = scenario("ListPositionSimulation")
    .exec(
      http("Aircraft list request")
        .get("/airspace/aircraft/list")
        .check(jsonPath("$.list").saveAs("icao24"))
    )
    .exec(
      http("Position endpoint url request")
        .get("/endpoints/position/url")
        .check(bodyString.saveAs("url"))
    )
    .exec(
      ws("Connect WS").connect("#{url}/airspace/aircraft/#{icao24(0)}/position").await(WS_CONNECTION_TIME)()
    )
    .pause(EXPLORE_API_USER_TIME)
    .exec(
      ws("Close WS").close
    )
    .pause(5)

  setUp(
    scn.inject(atOnceUsers(N_USERS))
  ).protocols(httpProtocol)
}
