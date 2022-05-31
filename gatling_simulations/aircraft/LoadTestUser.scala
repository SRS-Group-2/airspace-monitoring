package aircraft

import scala.concurrent.duration._

import io.gatling.core.Predef._
import io.gatling.http.Predef._
import io.gatling.http.check.ws.WsFrameCheck
class LoadTestUser extends Simulation {
  val BASE_URL = System.getProperty("base_url")
  val WS_URL = "wss://aircraft-positions-onwhlgspyq-ew.a.run.app"
  val N_USERS = Integer.getInteger("users", 20)
  val RAMP = java.lang.Long.getLong("ramp", 0)
  val WS_CONNECTION_TIME = 20
  val EXPLORE_API_USER_TIME = 45
  val nowTime: Long = System.currentTimeMillis / 1000
  val threeDaysAgoTime: Long = nowTime - 259200
  //val wsCheck = ws.checkTextMessage("checkName")

  val httpProtocol = http
    .baseUrl(BASE_URL)
    .wsBaseUrl(WS_URL)
    .acceptHeader("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
    .doNotTrackHeader("1")
    .acceptLanguageHeader("en-US,en;q=0.5")
    .acceptEncodingHeader("gzip, deflate")
    .userAgentHeader("Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")

  val scn = scenario("LoadTest")

    // web ui
    .exec(
      http("Web ui: airspace.css")
        .get("/airspace.css")
    )
    .exec(
      http("Web ui: index.js")
        .get("/index.js")
    )
    .exec(
      http("Web ui: favicon.ico")
        .get("/favicon.ico")
    )
    .exec(
      http("Web ui: hero-bg.png")
        .get("/hero-bg.png")
    )
    .exec(
      http("Web ui: plane.png")
        .get("/plane.png")
    )
    // request for aicraft_list
    .exec(
      http("Aircraft list request")
        .get("/airspace/aircraft/list")
        .check(jmesPath("icao24").transform(string => string.drop(1).drop(1).split("\",\"")).saveAs("icaoList"))
    )
    .pause(5) 
    // request for aircraft_info
    .exec(
      http("Aircraft info request")
        .get(s"/airspace/aircraft/#{icaoList.random()}/info")
    )
    .pause(5)

    // request for aicraft_daily_history (1h, 6h, 24h)
    .exec(
      http("Airspace recent history request (1h)")
        .get(s"/airspace/history/realtime/1h")
    )
    .pause(5)
    .exec(
      http("Airspace recent history request (6h)")
        .get(s"/airspace/history/realtime/6h")
    )
    .pause(5)
    .exec(
      http("Airspace recent history request (24h)")
        .get(s"/airspace/history/realtime/24h")
    )
    .pause(5)

    // request for aircraft_mothly_history (one with hour resolution, one with day resolution)
    .exec(
      http("Airspace old history request (resolution = day)")
        .get("/airspace/history?from=" + threeDaysAgoTime + "&to=" + nowTime + "&resolution=day")
    )
    .pause(5)
    .exec(
      http("Airspace old history request (resolution = hour)")
        .get("/airspace/history?from=" + threeDaysAgoTime + "&to=" + nowTime + "&resolution=hour")
    )
    .pause(5)

    //  request for aircraft_positions (wait 1.30 minutes for reading at least 6 diffent positions, then with another icao24 user does the same process)
    .exec(
      http("Position endpoint url request")
        .get("/endpoints/position/url")
    )
    .exec(
      ws("Connect WS").connect("/airspace/aircraft/#{icaoList.random()}/position")
    )
    .pause(EXPLORE_API_USER_TIME) //45
    .exec(
        ws("Close WS").close
    )
    .exec(
      ws("Connect WS").connect("/airspace/aircraft/#{icaoList.random()}/position")
    )
    .pause(30) // user iteraction 
    .exec(
        ws("Close WS").close
    )

  setUp(
    scn.inject(
      constantUsersPerSec(2).during(10.minutes)
    )
  ).protocols(httpProtocol)
}

