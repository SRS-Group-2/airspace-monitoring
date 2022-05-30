package aircraft

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._

class ListPositionSimulation extends Simulation {
  val BASE_URL = System.getProperty("base_url")
  val WS_URL = "wss://aircraft-positions-onwhlgspyq-ew.a.run.app"
  val N_USERS = Integer.getInteger("users", 20)
  val RAMP = java.lang.Long.getLong("ramp", 0)
  val WS_CONNECTION_TIME = 100
  val EXPLORE_API_USER_TIME = 60

  val httpProtocol = http
    .baseUrl(BASE_URL)
    .wsBaseUrl(WS_URL)
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .doNotTrackHeader("1")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")

  val scn = scenario("ListPositionSimulation")
    .exec(
      http("Aircraft list request")
        .get("/airspace/aircraft/list")
         .check(jmesPath("icao24").transform(string => string.drop(1).drop(1).split("\",\"")).saveAs("icaoList"))
    )
    .exec(
      http("Position endpoint url request for #{icaoList(0)}")
        .get("/endpoints/position/url")
    )
    .exec(
      ws("Connect WS").connect("/airspace/aircraft/#{icaoList(0)}/position")
    )
    .pause(EXPLORE_API_USER_TIME)
    .exec(
      ws("Close WS").close
    )

  setUp(
    scn.inject(
      atOnceUsers(3))
  ).protocols(httpProtocol)
}

