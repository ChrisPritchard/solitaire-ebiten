# Technical design

There are two overall approaches to card management:

- every card has an x,y and maybe z, and are dumbly rendered accordingly. game rules with stacks, tableaus etc is not represented in the state directly but instead in allowable actions and consequences. some additional state data might be required, like origin for moved cards
- cards are not objects at all, but instead each section of the board is represented as an array or nilable property. the only dynamic thing is anything being moved, or floated which might require additional data be kept like origin.

i think cards will either need to track what is stacked on them, or perhaps moveable stacks will need to exist as an object. this is to facilitate not just moving the bottom or free card of any stack, but in the case of the tableau piles, the ability to move the top most card where all stacked cards follow stacking rules.

to make the project more useful for other projects, i think embedding the rules directly into the structure is NOT appropriate. instead:

- track all cards by position and z-order
- each card also has a content which could be a go interface, with a draw function
- for solitaire, the content contains suit, value and whether its face up / shown. draw can then select the appropriate image

rules can be described as an object with a few methods maybe. board setup, can drag, drop action, has won etc. maybe a mvu system where events include game init, drag start, stop, click etc.

- init setups up the board, perhaps by first creating the stock and then emitting a series of move events
- when player clicks on a card, if the card can be moved its marked as moving in internal game state
- as the player drags, the moving card and all valid stacked are also moved
- if the player releases, check if it can be stacked in place. if not issue move orders back to position
- check victory conditions as appropriate

## dragging/mouse clicks

I want to support both mouse and touch inputs, which work differently in ebiten. To do this I need to track:

- whether these are pressed, or released
- the current state of a drag action
- if just pressed, is there a card that can be moved? if so, memorise it
- on updates, if the drag position has moved, then move the moving card
- on release, determine if the card can be placed there, if so place it, else snap back to its original position

suggestion: start with mouse and abstract to just whats needed: position, state (justpressed, dragging, released)
pass info to game rules, which can manage if something should occur

## initial setup

- shuffle deck
- render in position as a stack
- render spaces for foundation, waste etc
- deal out piles
  - move action moving one card at a time...
  -

## coordinate system

cards can only be in certain positions, when placed. they can be anyway while being dragged, but the placement positions on the board are fixed:

```
[]  []
[]  [] [] []
[]     [] []
[]        []
```

etc.

establishing this 'grid' would likely be a help

the grid is basically defined in card 'units', with whole units for things like the foundations, and partial units for the piles which overlay their contents

rather than using pixel coordinates, then, it might be advisable to store card positions in terms of units. however, how does this apply to dragging?

> is the separation between 'game' and the 'rules' getting in the way? or does it help keep things unique. i think i should drop it, for now, while still being mindful of separation. the main difference would be removing 'card content' from cards, and just embedding that data straight. also not relying entirely on a separate rules object for where things are placed. that being said, i like the moving of pure game stuff, like reading ebitengine's touch state, into its own thing. perhaps those could be removed into a helper class instead?
> there is also the separation of game objects from their image assets, which allows the 'game' to be more portable (if i do what i originally intended to do, and also build this in godot)
> and the fact that certain rendering considerations, like resolution, don't need to be a factor in the game. as a counter though, pixel positions are a mix, and sounds will also be a mix (though that might be solved with some sort of mvu eventing system)

a few modifications both to help with these grid mechanics, and dragging.

the idea is that a card tracks its board 'position', defined in card units. this could be, in the y direction, four units per card while in the x direction just two. the renderer part of things would handle converting this to pixel units
dragging becomes a little complex - while the game stores the drag state, and this could be modified to not use pixel coordinates on the card, how would that be surfaced to the renderer? presently cards are provided as a slice, so it would need to be (in the current structure) part of the card slice elements.

couple of options here; actually they're sort of the same option: use a different model for the presentation to the renderer. this would just be pixel positions and image asset - almost like a view model. this does take asset and size conversion out of the renderer and puts them in this view model transform, wherever that is, but that could be fine - it further simplifies the core game loop object.

## MVVM in go

how would this look?

we have the game mechanics, in one file
we have a view transform, that could just be a function
we have the renderer interfaces with ebitengine, including picking up the touch state on update

maybe keep Main, holding a GameLoop (instead of just Game, making somewhat clear it isnt the whole thing)
on update, the GameLoop would call a Transform function, passing in the Sawayama model object. This would no longer be a 'RuleSet' interface, maybe, but just a model or struct.

> in this way, the transform function doesnt get silly with its signature. it can be defined directly on sawayama. the close coupling might suggest they could be one object, but this distinctness allows for sawayama to be taken out and converted into another engine like the original concept.

the transform returns images, and maybe sounds? which the gameloop then draws.

dragging is interesting. the problem is the game returns pixel coordinates for the mouse, and whether or not the user is mouse clicking or touching. in the current model, the game uses those coords and its knowledge of both card size and positions to find an appropriate card to start dragging, at which point it tracks offsets etc. but if the new model doesnt contain most of this info in the rules, then how might this work?

the viewmodel could hold code to transform cursor positions into possible cards. so how would this work?

- the gameloop tracks pressed and cursor position
- the game really only needs to know whether a card has been selected or not, whether its being dragged. its position on the screen is irrelevant to the rules, until it is released at which point it needs to calculate if it can be placed in the position it is over
- the gameloop however, needs to render that image in a specific position. the viewmodel at present is what takes game rule cards and turns them into pixel coordinates
- the viewmodel is presently a function. could it be stateful? could it track the movement of a card? in a sense, as described above, the moving of cards around the screen is not a rule concern until they are released...

other considerations:

what about the deck? its a bunch of non-visible cards layered on top of each other. ideally this should look like a stack, and not just the last card's back. might require some smarts in the view model. could also just ignore this problem for now.

## remaining tasks for workable core rules

- click on deck to draw three to the side
- drag from waste
- drag onto foundations

the game rules themselves *could* split the cards into piles, foundations etc to have a little less iteration. but would that add much? I *do* need to understand if they have clicked on the deck, they're working off something in the waste (and therefore can only have the top card) or the foundation (and therefore can't drag from, but can drop to with the new cu coords being the same)

as a side note, right now touch integration is not passed to the game, only touch is detected. i guess there coulod be a 'is deck' event, or the viewmodel could understand where the deck is and click it.

## sorting issues

right now cards are rendered in the order they have within the deck. when cards are moved from a left pile to a right pile, their cu coords are updated but not their position in the deck, meaning that they get overdrawn by higher cards in the same pile

the solution could be:

- when moving, remove and insert into the card slice rather than (or in addition to) just updating cu
- have the deck resorted before rendering, or post stack move
- have the renderer (or the view model) follow rules to draw based on cux/cuy

none of these are super clean; the last one because it will require a sorting operation on every draw. or we could follow a 'dirty flag' model and do so only when the deck is changed. if the sorting is done via the viewmodel, it might be able to track this.

## empty spaces

when a pile is cleared, there is an empty space. also, when the deck is emptied. the draggable code needs to account for these somehow

i think the deck should be its own thing, too. perhaps revisit splitting cards into their positions? the difficulty is how we shift cards between piles, or other sets like foundations. the view model directly manipulating it seems unwise.

could change the dragging interaction, so instead of calling Droppable to see if cards can be dropped, just pass them to the game at the drop point and the game can handle updating them and such

## card movement and end game

So once you 'win', with your stacks all set, you need to currently manually move everything to the foundations which is painful. id like either a button or some game state detection that triggers an automatic stacking operation (and a final 'you win' banner)

> automatic card movement might also be used during dealing.

how would it work? I could just snap cards from deck to board, from board to foundations. this is what happens when the deck is clicked for the waste, right now. not super satisfying. how might an animation work?

the game might be frozen (no additional actions) while a card is moving. the view model can use its drag infrastructure to make the move. but if the game triggers it, and immediately updates, what does that do?

in terms of putting cards on the stack, the drag infra can be re-used. the view model would calculate what can be moved, initiate a drag state based on card positions, and initiate drops. it would need to use some sort of goroutine or stateful shift, though...

how would this work, using dragging? dragging tracks teh card, and its offset which is mouse position at drag start
instead, we could set drag state to the moved stack card, and set the offset to the cu-to-pixel conversion of that card.
we still need to track the destination cu, and the progress towards that destination. the progress might be used directly: calculate the dx, dy, apply the progress to those diffs, add to origin, and set card pixel position

# release ideas

pc, mac and android I think. lets see how easy this is

mobile is not easy, it seems. can build the game into a shared library, but then you need the native tooling to package them  (e.g. into a apk for android). something to investigate, but maybe... not now? if i learn how to do it now, just as i have learned out to do it in the past, the knowledge will no doubt be obsolete by the time i need to do it again.

pc and mac should be trivial though, and publisheable on github
