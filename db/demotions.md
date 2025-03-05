# Member Ranking Mechanisms #

As an alliance, we have had many different systems and objectives over the year
of playing this game that have determined member rankings, specifically
pertained to R1-R3.  In the game, there is currentky no functional difference
between those rankings other than pecieved value or recognition.  As the game
has matured we have seen that there are many different levels and motives to
play this game.

It is understood that every member is unique, and has different reasons for why
they play the game and why they are in P4K.  This mixture of motivations is a
benifit to the team for a variety of reasons.  However, one member's playing
style may appear to clash with other members' goals and this is creating a
divide within the alliance that is likely going to effect everyone's enjoyment.

For season 3 we have clear goals that have been set for the alliance as whole
from start.  In order to acheive those goals some members will need to adjust
their game play to ensure that those goals are met and to support the other
members that are sacrificing a great deal to achieve those goals.

To help define what it means to be a P4K member and to support as many member
objectives as possible, we will be implementing a new "R" system that will
hopefully clear up ambiguity around the most importnant things our members can
do to support one another and the alliance goals the allince is trying to
achieve.

## Ranking Definitions ##

Each rank from R3 to R1 will have a thresholds that dertmine which members
belong to each rank.  Each rank will also come with benefits or restritions for
which alliace events and activities can be participated in.  Lower ranked
players will be given less priority for rewards and other limited events, so it
will be important for each member to try their best to maintain an R3 ranking.

### R3 ###

All current members not currently at R4 will start at a ranking of R3.  A
ranking of R3 mean that all alliance standards are being met satisfactorally,
and should generally mean the member is contributing at high level.  Members at
the rank of R3 have earning the right to participate in all alliance activities,
and subject to supply will get all alliance rewards.

### R2 ###

Members at rank R2 have not met all of the P4K standards and thus have been
demoted from R3.  R2 status will mean less priority for alliance rewards and
limited activities like Desert Storm or Season rewards.  More importantly,
being ranked R2 means there are aspects of game play that the member needs to
improve on.

### R1 ###

A ranking or R1 indicates that the member has significantly failed to meet
alliance expectations and improvment is needed to continue membership.  R1
members will not be allowed to participate in any limited events and will get
the lowest alliance rewards.  Additionally, R1 members will also be asked to not
join in radar digs (except their own) or claim any alliance parties or eggs.

### R4 ###

Members at rank R4 are trusted with many responsabilities and are required to
maintain at least the minimum thresholds for R3.  R4 members will have similar,
but slightly higher, priority for limited alliance rewards and events than R3
members to compensate for the extra time and effort required for managing key
aspects of the alliance.

## Ranking Points and Thresholds ##

To determine a member's ranking, leadership will be tracking each members
contributions with regards to alliance participation and individual growth.
There will be a list of actions or inactions that a member can take to increase
or decrease the number of _ranking points_ they have.  Each of these actions will
have a value associated with them that will increase or decrease a member's
ranking score.  If a member's ranking score is below the threshold for their
current ranking on Sunday, that member will be demoted.

| Ranking | Points Threshold |
| -- | -- |
| R3 | 10 points |
| R2 | 5 points |
| R1 | 0 points |

Members below zero points will be removed from the alliance on Sunday, barring
extenuating cirumstances.

## Deeds and Offenses (Activities) ##

Not everything in Last War can be efficiently tracked, but leadership is
commited to tracking the participation metrics and behaviors that are possible
within reason.  Using this informatio, members will earn or lose rankning points
towards their total to detrmine what their overall rankning should be.  Each
activity that is tracked will be given a value to detrmine how many points can
be earned or lost.

### Activity Values ###

| Activity | Value |
| -- | -- |
| Top 10 overall in Alliance Duel (day) | 1 point |
| Top 3 overall in Alliance Duel (day) | 2 points |
| Alliance Duel MVP (day) | 3 points |
| Alliance Duel MVP (day+warzone) | 4 points |
| Arms race 1st place (warzone only) | 2 points |
| Logging in (day) | 1 point |
| Alliance donations over 60000 (week) | 1 point |
| Alliance donations under 30000 (week) | -1 point |
| More than 24 hours without logging in | -2 points |
| Less than 7.2 million in Alliance Duel (day) | -1 points |
| Less than 45 million in Alliance Duel (week) | -2 points |
| Not shielding and offline on Alliance Duel Buster Day | -5 points |
| Starting non-related rallies during alliance events | -1 point |
| Starting digs/eggs/parties during alliance events | -1 point |
| Mis-alignment in hive | -1 point |
| Ignoring alliance announcments or mail | -1 point |
| Ignoring urgent alliance announcements | -2 points |
| Sending blue secret tasks | -1 point |
| Sending non-gold secret tasks on buuld/kill day | -1 point |
| Disrespectful behavior towards alliance members | -5 points |
| Disrespectful behavior in general | -2 points |
| Other egregious conduct | -17 points |

## Ranking Objects ##

For technical use only :)

### Ranking Activities ###

A ranking activity is the definition of a transgression or deed, and includes
the following attributes:

* Description: the description of the activity
* Value: the amount of points the

### Ranking Activity Data ###

Activity data is a list storing the the activities performed, which
commander performed the activity, and the date the activity was performed:

* Member: the member who performed the activity
* Activity: the activity performed
* Date: the date of the activity

### Ranking Thresholds ###

Offense demotion stores the minimum number of offense points a commander needs
to earn to be placed at a certain rank.  The ranks are R3, R2, R1, and R0 (R0
indicates removal from the alliance)

* Rank: the ranking level threshold being defined
* Points: the minimum number of points to be placed at the given rank

## Conclusions ##

Using these objects, queries can be created to provide a record of a member's
offenses, deeds, the severity of those activities, and what ranking they should
currently hold in the alliance (or if they should be removed).
